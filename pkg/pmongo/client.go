package pmongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/lemonkingstar/spider/pkg/predis"
)

type Client struct {
	dbc    *mongo.Client
	dbname string
	sess   mongo.Session
	tm     *TxnManager
}

type MongoConf struct {
	TimeoutSeconds int
	MaxOpenConns   uint64
	MaxIdleConns   uint64
	URI            string
	RsName         string
	SocketTimeout  int
}

// NewMgo returns new RDB
func NewMgo(config MongoConf, timeout time.Duration) (*Client, error) {
	connStr, err := connstring.Parse(config.URI)
	if nil != err {
		return nil, err
	}
	if config.RsName == "" {
		return nil, fmt.Errorf("mongodb rsName not set")
	}
	socketTimeout := time.Second * time.Duration(config.SocketTimeout)
	// do not change this, our transaction plan need it to false.
	// it's related with the transaction number(eg txnNumber) in a transaction session.
	disableWriteRetry := false
	conOpt := options.ClientOptions{
		MaxPoolSize:    &config.MaxOpenConns,
		MinPoolSize:    &config.MaxIdleConns,
		ConnectTimeout: &timeout,
		SocketTimeout:  &socketTimeout,
		ReplicaSet:     &config.RsName,
		RetryWrites:    &disableWriteRetry,
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI), &conOpt)
	if nil != err {
		return nil, err
	}

	if err := client.Connect(context.TODO()); nil != err {
		return nil, err
	}

	// TODO: add this check later, this command needs authorize to get version.
	// if err := checkMongodbVersion(connStr.Database, client); err != nil {
	// 	return nil, err
	// }

	// initialize mongodb related metrics
	//initMongoMetric()

	return &Client{
		dbc:    client,
		dbname: connStr.Database,
		tm:     &TxnManager{},
	}, nil
}

// from now on, mongodb version must >= 4.2.0
func checkMongodbVersion(db string, client *mongo.Client) error {
	serverStatus, err := client.Database(db).RunCommand(
		context.Background(),
		bsonx.Doc{{"serverStatus", bsonx.Int32(1)}},
	).DecodeBytes()
	if err != nil {
		return err
	}

	version, err := serverStatus.LookupErr("version")
	if err != nil {
		return err
	}

	fields := strings.Split(version.StringValue(), ".")
	if len(fields) != 3 {
		return fmt.Errorf("got invalid mongodb version: %v", version.StringValue())
	}
	// version must be >= v4.2.0
	major, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s major failed, err: %v", version.StringValue(), err)
	}
	if major < 4 {
		return errors.New("mongodb version must be >= v4.2.0")
	}

	minor, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("parse mongodb version %s minor failed, err: %v", version.StringValue(), err)
	}
	if minor < 2 {
		return errors.New("mongodb version must be >= v4.2.0")
	}
	return nil
}

// InitTxnManager TxnID management of initial transaction
func (c *Client) InitTxnManager(r predis.Client) error {
	return c.tm.InitTxnManager(r)
}

// Close replica client
func (c *Client) Close() error {
	c.dbc.Disconnect(context.TODO())
	return nil
}

// Ping replica client
func (c *Client) Ping() error {
	return c.dbc.Ping(context.TODO(), nil)
}

// IsDuplicatedError check duplicated error
func (c *Client) IsDuplicatedError(err error) bool {
	if err != nil {
		if strings.Contains(err.Error(), "The existing index") {
			return true
		}
		if strings.Contains(err.Error(), "There's already an index with name") {
			return true
		}
		if strings.Contains(err.Error(), "E11000 duplicate") {
			return true
		}
		if strings.Contains(err.Error(), "IndexOptionsConflict") {
			return true
		}
		if strings.Contains(err.Error(), "all indexes already exist") {
			return true
		}
		if strings.Contains(err.Error(), "already exists with a different name") {
			return true
		}
	}
	return err == ErrDuplicated
}

// IsNotFoundError check the not found error
func (c *Client) IsNotFoundError(err error) bool {
	return err == ErrDocumentNotFound
}

// Table collection operation
func (c *Client) Table(collName string) Table {
	col := Collection{}
	col.collName = collName
	col.Client = c
	return &col
}

// get db client
func (c *Client) GetDBClient() *mongo.Client {
	return c.dbc
}

// get db name
func (c *Client) GetDBName() string {
	return c.dbname
}

// Collection implement client.Collection interface
type Collection struct {
	collName string // 集合名
	*Client
}

// Find 查询多个并反序列化到 Result
func (c *Collection) Find(filter Filter, opts ...FindOpts) Finder {
	find := &Find{
		Collection: c,
		filter:     filter,
		projection: make(map[string]int),
	}

	if len(opts) == 0 {
		find.projection["_id"] = 0
		return find
	}

	if !opts[0].WithObjectID {
		find.projection["_id"] = 0
		return find
	}
	return find
}

// Find define a find operation
type Find struct {
	*Collection

	projection map[string]int
	filter     Filter
	start      int64
	limit      int64
	sort       bson.D
}

// Fields 查询字段
func (f *Find) Fields(fields ...string) Finder {
	for _, field := range fields {
		if len(field) <= 0 {
			continue
		}
		f.projection[field] = 1
	}
	return f
}

// Sort 查询排序
// sort支持多字段最左原则排序
// sort值为"host_id, -host_name"和sort值为"host_id:1, host_name:-1"是一样的，都代表先按host_id递增排序，再按host_name递减排序
func (f *Find) Sort(sort string) Finder {
	if sort != "" {
		sortArr := strings.Split(sort, ",")
		f.sort = bson.D{}
		for _, sortItem := range sortArr {
			sortItemArr := strings.Split(strings.TrimSpace(sortItem), ":")
			sortKey := strings.TrimLeft(sortItemArr[0], "+-")
			if len(sortItemArr) == 2 {
				sortDescFlag := strings.TrimSpace(sortItemArr[1])
				if sortDescFlag == "-1" {
					f.sort = append(f.sort, bson.E{sortKey, -1})
				} else {
					f.sort = append(f.sort, bson.E{sortKey, 1})
				}
			} else {
				if strings.HasPrefix(sortItemArr[0], "-") {
					f.sort = append(f.sort, bson.E{sortKey, -1})
				} else {
					f.sort = append(f.sort, bson.E{sortKey, 1})
				}
			}
		}
	}

	return f
}

// Start 查询上标
func (f *Find) Start(start uint64) Finder {
	// change to int64,后续改成int64
	dbStart := int64(start)
	f.start = dbStart
	return f
}

// Limit 查询限制
func (f *Find) Limit(limit uint64) Finder {
	// change to int64,后续改成int64
	dbLimit := int64(limit)
	f.limit = dbLimit
	return f
}

// All 查询多个
func (f *Find) All(ctx context.Context, result interface{}) error {
	findOpts := &options.FindOptions{}
	if len(f.projection) != 0 {
		findOpts.Projection = f.projection
	}
	if f.start != 0 {
		findOpts.SetSkip(f.start)
	}
	if f.limit != 0 {
		findOpts.SetLimit(f.limit)
	}
	if len(f.sort) != 0 {
		findOpts.SetSort(f.sort)
	}
	// 查询条件为空时候，mongodb 不返回数据
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	return f.tm.AutoRunWithTxn(ctx, f.dbc, func(ctx context.Context) error {
		cursor, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			return err
		}
		return cursor.All(ctx, result)
	})
}

// One 查询一个
func (f *Find) One(ctx context.Context, result interface{}) error {
	findOpts := &options.FindOptions{}
	if len(f.projection) != 0 {
		findOpts.Projection = f.projection
	}
	if f.start != 0 {
		findOpts.SetSkip(f.start)
	}
	if f.limit != 0 {
		findOpts.SetLimit(1)
	}
	if len(f.sort) != 0 {
		findOpts.SetSort(f.sort)
	}
	// 查询条件为空时候，mongodb panic
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)
	return f.tm.AutoRunWithTxn(ctx, f.dbc, func(ctx context.Context) error {
		cursor, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).Find(ctx, f.filter, findOpts)
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			return cursor.Decode(result)
		}
		return ErrDocumentNotFound
	})

}

// Count 统计数量(非事务)
func (f *Find) Count(ctx context.Context) (uint64, error) {
	if f.filter == nil {
		f.filter = bson.M{}
	}

	opt := getCollectionOption(ctx)

	sessCtx, _, useTxn, err := f.tm.GetTxnContext(ctx, f.dbc)
	if err != nil {
		return 0, err
	}
	if !useTxn {
		// not use transaction.
		cnt, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).CountDocuments(ctx, f.filter)
		if err != nil {
			return 0, err
		}

		return uint64(cnt), err
	} else {
		// use transaction
		cnt, err := f.dbc.Database(f.dbname).Collection(f.collName, opt).CountDocuments(sessCtx, f.filter)
		// do not release th session, otherwise, the session will be returned to the
		// session pool and will be reused. then mongodb driver will increase the transaction number
		// automatically and do read/write retry if policy is set.
		// mongo.CmdbReleaseSession(ctx, session)
		if err != nil {
			return 0, err
		}
		return uint64(cnt), nil
	}
}

// Insert 插入数据, docs 可以为 单个数据 或者 多个数据
func (c *Collection) Insert(ctx context.Context, docs interface{}) error {

	rows := ConvertToInterfaceSlice(docs)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).InsertMany(ctx, rows)
		if err != nil {
			return err
		}

		return nil
	})
}

// Update 更新数据
func (c *Collection) Update(ctx context.Context, filter Filter, doc interface{}) error {
	if filter == nil {
		filter = bson.M{}
	}

	data := bson.M{"$set": doc}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, data)
		if err != nil {
			return err
		}
		return nil
	})
}

// Upsert 数据存在更新数据，否则新加数据。
// 注意：该接口非原子操作，可能存在插入多条相同数据的风险。
func (c *Collection) Upsert(ctx context.Context, filter Filter, doc interface{}) error {
	// set upsert option
	doUpsert := true
	replaceOpt := &options.UpdateOptions{
		Upsert: &doUpsert,
	}
	data := bson.M{"$set": doc}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateOne(ctx, filter, data, replaceOpt)
		if err != nil {
			return err
		}
		return nil
	})

}

// UpdateMultiModel 根据不同的操作符去更新数据
func (c *Collection) UpdateMultiModel(ctx context.Context, filter Filter, updateModel ...ModeUpdate) error {
	data := bson.M{}
	for _, item := range updateModel {
		if _, ok := data[item.Op]; ok {
			return errors.New(item.Op + " appear multiple times")
		}
		data["$"+item.Op] = item.Doc
	}

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, data)
		if err != nil {
			return err
		}
		return nil
	})

}

// Delete 删除数据
func (c *Collection) Delete(ctx context.Context, filter Filter) error {
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).DeleteMany(ctx, filter)
		if err != nil {
			return err
		}

		return nil
	})

}

type Idgen struct {
	ID         string `bson:"_id"`
	SequenceID uint64 `bson:"SequenceID"`
}

// HasTable 判断是否存在集合
func (c *Client) HasTable(ctx context.Context, collName string) (bool, error) {
	cursor, err := c.dbc.Database(c.dbname).ListCollections(ctx, bson.M{"name": collName, "type": "collection"})
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		return true, nil
	}

	return false, nil
}

// DropTable 移除集合
func (c *Client) DropTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).Collection(collName).Drop(ctx)
}

// CreateTable 创建集合 TODO test
func (c *Client) CreateTable(ctx context.Context, collName string) error {
	return c.dbc.Database(c.dbname).RunCommand(ctx, map[string]interface{}{"create": collName}).Err()
}

// CreateIndex 创建索引
func (c *Collection) CreateIndex(ctx context.Context, index Index) error {
	createIndexOpt := &options.IndexOptions{
		Background: &index.Background,
		Unique:     &index.Unique,
	}
	if index.Name != "" {
		createIndexOpt.Name = &index.Name
	}

	if index.ExpireAfterSeconds != 0 {
		createIndexOpt.SetExpireAfterSeconds(index.ExpireAfterSeconds)
	}

	createIndexInfo := mongo.IndexModel{
		Keys:    index.Keys,
		Options: createIndexOpt,
	}

	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.CreateOne(ctx, createIndexInfo)
	if err != nil {
		// ignore the following case
		// 1.the new index is exactly the same as the existing one
		// 2.the new index has same keys with the existing one, but its name is different
		if strings.Contains(err.Error(), "all indexes already exist") ||
			strings.Contains(err.Error(), "already exists with a different name") {
			return nil
		}
	}

	return err
}

// DropIndex remove index by name
func (c *Collection) DropIndex(ctx context.Context, indexName string) error {
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	_, err := indexView.DropOne(ctx, indexName)
	return err
}

// Indexes get all indexes for the collection
func (c *Collection) Indexes(ctx context.Context) ([]Index, error) {
	indexView := c.dbc.Database(c.dbname).Collection(c.collName).Indexes()
	cursor, err := indexView.List(ctx)
	if nil != err {
		return nil, err
	}
	defer cursor.Close(ctx)
	var indexs []Index
	for cursor.Next(ctx) {
		idxResult := Index{}
		cursor.Decode(&idxResult)
		indexs = append(indexs, idxResult)
	}

	return indexs, nil
}

// AddColumn add a new column for the collection
func (c *Collection) AddColumn(ctx context.Context, column string, value interface{}) error {
	selector := Document{column: Document{"$exists": false}}
	datac := Document{"$set": Document{column: value}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, selector, datac)
		if err != nil {
			return err
		}
		return nil
	})
}

// RenameColumn rename a column for the collection
func (c *Collection) RenameColumn(ctx context.Context, filter Filter, oldName, newColumn string) error {
	if filter == nil {
		filter = Document{}
	}

	datac := Document{"$rename": Document{oldName: newColumn}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			return err
		}

		return nil
	})
}

// DropColumn remove a column by the name
func (c *Collection) DropColumn(ctx context.Context, field string) error {
	datac := Document{"$unset": Document{field: ""}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, Document{}, datac)
		if err != nil {
			return err
		}

		return nil
	})
}

// DropColumns remove many columns by the name
func (c *Collection) DropColumns(ctx context.Context, filter Filter, fields []string) error {

	unsetFields := make(map[string]interface{})
	for _, field := range fields {
		unsetFields[field] = ""
	}

	datac := Document{"$unset": unsetFields}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			return err
		}

		return nil
	})
}

// DropDocsColumn remove a column by the name for doc use filter
func (c *Collection) DropDocsColumn(ctx context.Context, field string, filter Filter) error {
	// 查询条件为空时候，mongodb 不返回数据
	if filter == nil {
		filter = bson.M{}
	}

	datac := Document{"$unset": Document{field: ""}}
	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		_, err := c.dbc.Database(c.dbname).Collection(c.collName).UpdateMany(ctx, filter, datac)
		if err != nil {
			return err
		}

		return nil
	})
}

// AggregateAll aggregate all operation
func (c *Collection) AggregateAll(ctx context.Context, pipeline interface{}, result interface{}) error {
	opt := getCollectionOption(ctx)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		cursor, err := c.dbc.Database(c.dbname).Collection(c.collName, opt).Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)
		return decodeCusorIntoSlice(ctx, cursor, result)
	})

}

// AggregateOne aggregate one operation
func (c *Collection) AggregateOne(ctx context.Context, pipeline interface{}, result interface{}) error {
	opt := getCollectionOption(ctx)

	return c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		cursor, err := c.dbc.Database(c.dbname).Collection(c.collName, opt).Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		defer cursor.Close(ctx)
		for cursor.Next(ctx) {
			return cursor.Decode(result)
		}
		return ErrDocumentNotFound
	})

}

// Distinct Finds the distinct values for a specified field across a single collection or view and returns the results in an
// field the field for which to return distinct values.
// filter query that specifies the documents from which to retrieve the distinct values.
func (c *Collection) Distinct(ctx context.Context, field string, filter Filter) ([]interface{}, error) {
	if filter == nil {
		filter = bson.M{}
	}

	opt := getCollectionOption(ctx)
	var results []interface{} = nil
	err := c.tm.AutoRunWithTxn(ctx, c.dbc, func(ctx context.Context) error {
		var err error
		results, err = c.dbc.Database(c.dbname).Collection(c.collName, opt).Distinct(ctx, field, filter)
		if err != nil {
			return err
		}

		return nil
	})
	return results, err
}

func decodeCusorIntoSlice(ctx context.Context, cursor *mongo.Cursor, result interface{}) error {
	resultv := reflect.ValueOf(result)
	if resultv.Kind() != reflect.Ptr || resultv.Elem().Kind() != reflect.Slice {
		return errors.New("result argument must be a slice address")
	}

	elemt := resultv.Elem().Type().Elem()
	slice := reflect.MakeSlice(resultv.Elem().Type(), 0, 10)
	for cursor.Next(ctx) {
		elemp := reflect.New(elemt)
		if err := cursor.Decode(elemp.Interface()); nil != err {
			return err
		}
		slice = reflect.Append(slice, elemp.Elem())
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	resultv.Elem().Set(slice)
	return nil
}

const (
	// reference doc:
	// https://docs.mongodb.com/manual/core/read-preference-staleness/#replica-set-read-preference-max-staleness
	// this is the minimum value of maxStalenessSeconds allowed.
	// specifying a smaller maxStalenessSeconds value will raise an error. Clients estimate secondaries’ staleness
	// by periodically checking the latest write date of each replica set member. Since these checks are infrequent,
	// the staleness estimate is coarse. Thus, clients cannot enforce a maxStalenessSeconds value of less than
	// 90 seconds.
	maxStalenessSeconds = 90 * time.Second
)

func getCollectionOption(ctx context.Context) *options.CollectionOptions {
	var opt *options.CollectionOptions
	switch GetDBReadPreference(ctx) {

	case NilMode:

	case PrimaryMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Primary(),
		}
	case PrimaryPreferredMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.PrimaryPreferred(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case SecondaryMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Secondary(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case SecondaryPreferredMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.SecondaryPreferred(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	case NearestMode:
		opt = &options.CollectionOptions{
			ReadPreference: readpref.Nearest(readpref.WithMaxStaleness(maxStalenessSeconds)),
		}
	}

	return opt
}
