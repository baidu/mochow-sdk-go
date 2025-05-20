package main

import (
	"fmt"
	"log"
	"time"

	"github.com/baidu/mochow-sdk-go/v2/client"
	"github.com/baidu/mochow-sdk-go/v2/mochow"
	"github.com/baidu/mochow-sdk-go/v2/mochow/api"
)

type MochowTest struct {
	database string
	table    string
	client   *mochow.Client
}

func NewMochowTest(database, table string, config *mochow.ClientConfiguration) (*MochowTest, error) {
	client, err := mochow.NewClientWithConfig(config)
	if err != nil {
		return nil, err
	}
	mochowTest := &MochowTest{database: database, table: table, client: client}
	return mochowTest, nil
}

func (m *MochowTest) clearEnv() error {
	// skip drop database and table when not existed
	dbExisted, err := m.client.HasDatabase(m.database)
	if err != nil {
		log.Fatalf("Fail to check database due to error: %v", err)
		return err
	}
	if !dbExisted {
		return nil
	}

	// drop table first if existed
	tableExisted, err := m.client.HasTable(m.database, m.table)
	if err != nil {
		log.Fatalf("Fail to check table due to error: %v", err)
		return err
	}
	if tableExisted {
		if err := m.client.DropTable(m.database, m.table); err != nil {
			log.Printf("Fail to drop table due to error: %v", err)
			return err
		}
		// wait util drop finished
		for {
			time.Sleep(time.Second * 5)
			if _, err := m.client.DescTable(m.database, m.table); err != nil {
				if realErr, ok := err.(*client.BceServiceError); ok {
					if realErr.StatusCode == 404 || realErr.Code == int(api.TableNotExist) {
						log.Println("Table already dropped")
						break
					}
				}
			}
		}
	}

	// drop database if existed
	if dbExisted {
		log.Printf("Try to drop existed database: %s", m.database)
		if err := m.client.DropDatabase(m.database); err != nil {
			log.Fatalf("Fail to drop database due to error: %v", err)
			return err
		}
	}

	_ = m.client.DropUser("user1")

	_ = m.client.DropRole("writable")

	return nil
}

func (m *MochowTest) createDatabaseAndTable() error {
	// create database
	if err := m.client.CreateDatabase(m.database); err != nil {
		log.Fatalf("Fail to create database due to error: %v", err)
		return err
	}

	// Fields
	fields := []api.FieldSchema{
		{
			FieldName:     "id",
			FieldType:     api.FieldTypeString,
			PrimaryKey:    true,
			PartitionKey:  true,
			AutoIncrement: false,
			NotNull:       true,
		},
		{
			FieldName: "bookName",
			FieldType: api.FieldTypeString,
			NotNull:   true,
		},
		{
			FieldName: "author",
			FieldType: api.FieldTypeString,
			NotNull:   true,
		},
		{
			FieldName: "page",
			FieldType: api.FieldTypeUint32,
			NotNull:   true,
		},
		{
			FieldName: "vector",
			FieldType: api.FieldTypeFloatVector,
			NotNull:   true,
			Dimension: 3,
		},
		{
			FieldName: "segment",
			FieldType: api.FieldTypeText,
			NotNull:   true,
		},
		{
			FieldName:   "arrField",
			FieldType:   api.FieldTypeArray,
			ElementType: api.ElementTypeString,
			NotNull:     true,
		},
	}

	// Indexes
	autoBuildPolicy := api.NewAutoBuildIncrementPolicy()
	autoBuildPolicy.AddRowCountIncrement(5000)
	indexes := []api.IndexSchema{
		{
			IndexName: "book_name_idx",
			Field:     "bookName",
			IndexType: api.SecondaryIndex,
		},
		{
			IndexName: "compound_filtering_idx",
			IndexType: api.FilteringIndex,
			FilterIndexFields: []api.FilteringIndexField{
				{
					Field:              "bookName",
					IndexStructureType: api.IndexStructureTypeDefault,
				},
				{
					Field:              "page",
					IndexStructureType: api.IndexStructureTypeDefault,
				},
			},
		},
		{
			IndexName:  "vector_idx",
			Field:      "vector",
			IndexType:  api.HNSW,
			MetricType: api.L2,
			Params: api.VectorIndexParams{
				"M":              32,
				"efConstruction": 200,
			},
			AutoBuild:       true,
			AutoBuildPolicy: autoBuildPolicy.Params(),
		},
		{
			IndexName:                    "book_segment_inverted_idx",
			IndexType:                    api.InvertedIndex,
			InvertedIndexFields:          []string{"segment"},
			InvertedIndexFieldAttributes: []api.InvertedIndexFieldAttribute{api.Analyzed},
			Params: api.NewInvertedIndexParams().
				Analyzer(api.ChineseAnalyzer).
				ParseMode(api.FineMode).
				AnalyzerCaseSensitive(true),
		},
	}

	// create table
	createTableArgs := &api.CreateTableArgs{
		Database:    m.database,
		Table:       m.table,
		Description: "basic test",
		Replication: 3,
		Partition: &api.PartitionParams{
			PartitionType: api.HASH,
			PartitionNum:  3,
		},
		EnableDynamicField: false,
		Schema: &api.TableSchema{
			Fields:  fields,
			Indexes: indexes,
		},
	}
	if err := m.client.CreateTable(createTableArgs); err != nil {
		log.Fatalf("Fail to create table due to error: %v", err)
		return err
	}
	for {
		time.Sleep(5 * time.Second)
		describeTableResult, err := m.client.DescTable(m.database, m.table)
		if err == nil && describeTableResult.Table.State == api.TableStateNormal {
			log.Println("Table create finished")
			break
		}
	}
	return nil
}

func (m *MochowTest) upsertData() error {
	data := []api.Row{
		{
			Fields: map[string]interface{}{
				"id":       "0001",
				"bookName": "西游记",
				"author":   "吴承恩",
				"page":     21,
				"vector":   []float32{0.2123, 0.21, 0.213},
				"arrField": []string{},
				"segment":  "富贵功名，前缘分定，为人切莫欺心。",
			},
		},
		{
			Fields: map[string]interface{}{
				"id":       "0002",
				"bookName": "西游记",
				"author":   "吴承恩",
				"page":     22,
				"vector":   []float32{0.2123, 0.22, 0.213},
				"arrField": []string{},
				"segment":  "正大光明，忠良善果弥深。些些狂妄天加谴，眼前不遇待时临。",
			},
		},
		{
			Fields: map[string]interface{}{
				"id":       "0003",
				"bookName": "三国演义",
				"author":   "罗贯中",
				"page":     23,
				"vector":   []float32{0.2123, 0.23, 0.213},
				"arrField": []string{"细作", "吕布"},
				"segment":  "细作探知这个消息，飞报吕布。",
			},
		},
		{
			Fields: map[string]interface{}{
				"id":       "0004",
				"bookName": "三国演义",
				"author":   "罗贯中",
				"page":     24,
				"vector":   []float32{0.2123, 0.24, 0.213},
				"arrField": []string{"吕布", "陈宫", "刘玄德"},
				"segment":  "布大惊，与陈宫商议。宫曰：“闻刘玄德新领徐州，可往投之。” 布从其言，竟投徐州来。有人报知玄德。",
			},
		},
		{
			Fields: map[string]interface{}{
				"id":       "0005",
				"bookName": "三国演义",
				"author":   "罗贯中",
				"page":     25,
				"vector":   []float32{0.2123, 0.24, 0.213},
				"arrField": []string{"玄德", "吕布", "糜竺"},
				"segment":  "玄德曰：“布乃当今英勇之士，可出迎之。”糜竺曰：“吕布乃虎狼之徒，不可收留；收则伤人矣。",
			},
		},
	}
	upsertArgs := &api.UpsertRowArg{
		Database: m.database,
		Table:    m.table,
		Rows:     data,
	}
	upsertResult, err := m.client.UpsertRow(upsertArgs)
	if err != nil {
		log.Fatalf("Fail to upsert row due to error: %v", err)
		return err
	}
	log.Printf("Upsert row result: %+v", upsertResult)
	return nil
}

func (m *MochowTest) queryData() error {
	queryArgs := &api.QueryRowArgs{
		Database: m.database,
		Table:    m.table,
		PrimaryKey: map[string]interface{}{
			"id": "0001",
		},
		Projections:    []string{"id", "bookName"},
		RetrieveVector: false,
	}
	queryResult, err := m.client.QueryRow(queryArgs)
	if err != nil {
		log.Fatalf("Fail to query row due to error: %v", err)
		return err
	}
	log.Printf("Query row result: %+v", queryResult)
	return nil
}

func (m *MochowTest) batchQueryData() error {
	batchQueryArgs := &api.BatchQueryRowArgs{
		Database: m.database,
		Table:    m.table,
		Keys: []api.QueryKey{
			{
				PrimaryKey: map[string]interface{}{
					"id": "0001",
				},
			},
			{
				PrimaryKey: map[string]interface{}{
					"id": "0002",
				},
			},
		},
		Projections:    []string{"id", "bookName"},
		RetrieveVector: false,
	}
	queryResult, err := m.client.BatchQueryRow(batchQueryArgs)
	if err != nil {
		log.Fatalf("Fail to batch query row due to error: %v", err)
		return err
	}
	log.Printf("batch query row result: %+v", queryResult)
	return nil
}

func (m *MochowTest) selectData() error {
	selectArgs := &api.SelectRowArgs{
		Database:    m.database,
		Table:       m.table,
		Projections: []string{"id", "bookName"},
		Limit:       1,
	}
	for {
		selectResult, err := m.client.SelectRow(selectArgs)
		if err != nil {
			log.Fatalf("Fail to select row due to error: %v", err)
			return err
		}
		log.Printf("Select row result: %+v", selectResult)
		if !selectResult.IsTruncated {
			break
		}
		selectArgs.Marker = selectResult.NextMarker
	}
	return nil
}

func (m *MochowTest) updateData() error {
	updateArgs := &api.UpdateRowArgs{
		Database: m.database,
		Table:    m.table,
		PrimaryKey: map[string]interface{}{
			"id": "0001",
		},
		Update: map[string]interface{}{
			"bookName": "红楼梦",
			"author":   "曹雪芹",
			"page":     100,
			"segment":  "满纸荒唐言，一把辛酸泪",
		},
	}
	err := m.client.UpdateRow(updateArgs)
	if err != nil {
		log.Fatalf("Fail to update row due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) topkSearch() error {
	// rebuild vector index
	if err := m.client.RebuildIndex(m.database, m.table, "vector_idx"); err != nil {
		log.Fatalf("Fail to rebuild index due to error: %v", err)
		return err
	}
	for {
		time.Sleep(5 * time.Second)
		descIndexResult, _ := m.client.DescIndex(m.database, m.table, "vector_idx")
		if descIndexResult.Index.State == api.IndexStateNormal {
			log.Println("Index rebuild finished")
			break
		}
	}

	// search
	vector := api.FloatVector([]float32{0.3123, 0.43, 0.213})

	searchArgs := &api.VectorSearchArgs{
		Database: m.database,
		Table:    m.table,
		Request: api.VectorTopkSearchRequest{}.
			New("vector", vector, 5).
			Filter("bookName='三国演义'").
			Config(api.VectorSearchConfig{}.New().Ef(200)),
	}

	searchResult, err := m.client.VectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search row due to error: %v", err)
		return err
	}

	log.Printf("topk search result: %+v", searchResult.Rows)
	return nil
}

func (m *MochowTest) rangeSearch() error {
	vector := api.FloatVector([]float32{0.3123, 0.43, 0.213})

	searchArgs := &api.VectorSearchArgs{
		Database: m.database,
		Table:    m.table,
		Request: api.VectorRangeSearchRequest{}.
			New("vector", vector, api.DistanceRange{Min: 0, Max: 20}).
			Filter("bookName='三国演义'").
			Limit(15).
			Config(api.VectorSearchConfig{}.New().Ef(200)),
	}

	searchResult, err := m.client.VectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search row due to error: %v", err)
		return err
	}

	log.Printf("range search result: %+v", searchResult.Rows)
	return nil
}

func (m *MochowTest) batchSearch() error {
	vectors := []api.Vector{
		api.FloatVector{0.3123, 0.43, 0.213},
		api.FloatVector{0.5512, 0.33, 0.43},
	}

	searchArgs := &api.VectorSearchArgs{
		Database: m.database,
		Table:    m.table,
		Request: api.VectorBatchSearchRequest{}.New("vector", vectors).
			Filter("bookName='三国演义'").
			Limit(10).
			Config(api.VectorSearchConfig{}.New().Ef(200)).
			Projections([]string{"id", "bookName", "author", "page"}),
	}

	searchResult, err := m.client.VectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to batch search row due to error: %v", err)
		return err
	}

	log.Printf("batch search result: %+v", searchResult.BatchRows)

	return nil
}

func (m *MochowTest) bm25Search() error {
	searchArgs := &api.BM25SearchArgs{
		Database: m.database,
		Table:    m.table,
		Request: api.BM25SearchRequest{}.
			New("book_segment_inverted_idx", "吕布").
			Limit(12).
			Filter("bookName='三国演义'").
			ReadConsistency(api.ReadConsistencyStrong).
			Projections([]string{"id", "vector"}),
	}

	searchResult, err := m.client.BM25Search(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search row due to error: %v", err)
		return err
	}

	log.Printf("bm25 search result: %+v", searchResult.Rows)
	return nil
}

func (m *MochowTest) hybridSearch() error {
	vector := api.FloatVector([]float32{0.3123, 0.43, 0.213})

	request := api.HybridSearchRequest{}.
		New(
			api.VectorTopkSearchRequest{}.New("vector", vector, 15).Config(api.VectorSearchConfig{}.New().Ef(200)), /*vector search args*/
			api.BM25SearchRequest{}.New("book_segment_inverted_idx", "吕布"),                                         /*BM25 search args*/
			0.4, /*vector search weight*/
			0.6 /*BM25 search weight*/).
		Filter("bookName='三国演义'").
		Limit(15)

	searchArgs := &api.HybridSearchArgs{
		Database: m.database,
		Table:    m.table,
		Request:  request,
	}

	searchResult, err := m.client.HybridSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search row due to error: %v", err)
		return err
	}

	log.Printf("hybrid search result: %+v", searchResult.Rows)
	return nil
}

func (m *MochowTest) multiVectorSearch() error {
	vector := api.FloatVector([]float32{0.3123, 0.43, 0.213})

	request := api.MultivectorSearchRequest{}.New().
		AddSingleVectorSearchRequest(api.VectorTopkSearchRequest{}.New("vector", vector, 15).Config(api.VectorSearchConfig{}.New().Ef(200))).
		AddSingleVectorSearchRequest(api.VectorTopkSearchRequest{}.New("vector_2", vector, 10).Config(api.VectorSearchConfig{}.New().Ef(200))).
		Rank(api.WeightedRank{}.New([]float64{0.123})).
		Filter("bookName='三国演义'").
		Limit(15)

	searchArgs := &api.MultivectorSearchArgs{
		Database: m.database,
		Table:    m.table,
		Request:  request,
	}

	searchResult, err := m.client.MultivectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search row due to error: %v", err)
		return err
	}

	log.Printf("hybrid search result: %+v", searchResult.Rows)
	return nil
}

func (m *MochowTest) searchIterator() error {
	// Example 1: TopK search iterator using SearchIteratorArgs
	log.Println("Start search iterator for TopK search with SearchIteratorArgs")
	vector := api.FloatVector([]float32{0.3123, 0.43, 0.213})

	batchSize := uint32(1000)
	request := api.VectorTopkSearchRequest{}.
		New("vector", vector, batchSize).
		Filter("bookName='三国演义'").
		Config(api.VectorSearchConfig{}.New().Ef(200))

	// Create iterator args
	iteratorArgs := &api.SearchIteratorArgs{
		Database:        m.database,
		Table:           m.table,
		Request:         request,
		BatchSize:       batchSize,
		TotalSize:       10000,
		Projections:     []string{"id", "bookName", "author"},
		ReadConsistency: api.ReadConsistencyEventual,
	}

	// Create iterator
	iterator, err := m.client.SearchIterator(iteratorArgs)
	if err != nil {
		log.Fatalf("Failed to create search iterator: %v", err)
		return err
	}

	// Iterate over search results
	totalRows := 0
	for {
		rows, err := iterator.Next()
		if err != nil {
			log.Fatalf("Failed to fetch search results: %v", err)
			return err
		}
		if len(rows) == 0 {
			break
		}
		log.Printf("Received batch with %d rows", len(rows))
		totalRows += len(rows)
	}
	log.Printf("Retrieved a total of %d rows", totalRows)
	iterator.Close()
	log.Println("Finished search iterator for TopK search")

	return nil
}

func (m *MochowTest) deleteDataWithPK() error {
	deleteArgs := &api.DeleteRowArgs{
		Database: m.database,
		Table:    m.table,
		PrimaryKey: map[string]interface{}{
			"id": "0001",
		},
	}
	if err := m.client.DeleteRow(deleteArgs); err != nil {
		log.Fatalf("Fail to delete row due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) deleteDataWithFilter() error {
	deleteArgs := &api.DeleteRowArgs{
		Database: m.database,
		Table:    m.table,
		Filter:   "id = '0002'",
	}
	if err := m.client.DeleteRow(deleteArgs); err != nil {
		log.Fatalf("Fail to delete row due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) dropAndCreateVIndex() error {
	// drop vindex
	if err := m.client.DropIndex(m.database, m.table, "vector_idx"); err != nil {
		log.Fatalf("Fail to drop index due to error: %v", err)
		return err
	}
	for {
		time.Sleep(5 * time.Second)
		_, err := m.client.DescIndex(m.database, m.table, "vector_idx")
		if realErr, ok := err.(*client.BceServiceError); ok {
			if realErr.Code == int(api.IndexNotExist) {
				log.Print("Index already dropped")
				break
			}
		}
	}

	// create vector index
	createIndexArgs := &api.CreateIndexArgs{
		Database: m.database,
		Table:    m.table,
		Indexes: []api.IndexSchema{
			{
				IndexName:  "vector_idx",
				Field:      "vector",
				IndexType:  api.HNSW,
				MetricType: api.L2,
				Params: api.VectorIndexParams{
					"M":              16,
					"efConstruction": 200,
				},
			},
		},
	}
	if err := m.client.CreateIndex(createIndexArgs); err != nil {
		log.Fatalf("Fail to create index due to error: %v", err)
		return err
	}
	return nil
}

// /////////////////////////////// rbac ///////////////////////////////////////
func (m *MochowTest) createRole() error {
	if err := m.client.CreateRole("writable"); err != nil {
		log.Fatalf("Fail to create role due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) grantRolePrivileges() error {

	args := &api.GrantRolePrivilegesArgs{
		Role: "writable",
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "*",
				Table:      "*",
				Privileges: []string{"QUERY", "SELECT", "SEARCH"},
			},
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_CONTROL", "TABLE_READWRITE"},
			},
		},
	}
	if err := m.client.GrantRolePrivileges(args); err != nil {
		log.Fatalf("Fail to GrantRolePrivileges due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) revokeRolePrivileges() error {

	args := &api.RevokeRolePrivilegesArgs{
		Role: "writable",
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_READWRITE"},
			},
		},
	}
	if err := m.client.RevokeRolePrivileges(args); err != nil {
		log.Fatalf("Fail to RevokeRolePrivileges due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) showRolePrivileges() error {

	ret, err := m.client.ShowRolePrivileges("writable")
	if err != nil {
		log.Fatalf("Fail to RevokeRolePrivileges due to error: %v", err)
		return err
	}
	fmt.Println("Privileges: ", ret)
	return nil
}

func (m *MochowTest) SelectRole() error {
	args := &api.SelectRoleArgs{
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_CONTROL"},
			},
		},
	}
	ret, err := m.client.SelectRole(args)
	if err != nil {
		log.Fatalf("Fail to SelectRole due to error: %v", err)
		return err
	}
	fmt.Println("Role: ", ret)
	return nil
}

func (m *MochowTest) dropRole() error {
	if err := m.client.DropRole("writable"); err != nil {
		log.Fatalf("Fail to drop user due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) createUser() error {
	if err := m.client.CreateUser("user1", "mochow@123"); err != nil {
		log.Fatalf("Fail to create user due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) grantUserRoles() error {
	args := &api.GrantUserRolesArgs{
		Username: "user1",
		Roles:    []string{"writable"},
	}
	if err := m.client.GrantUserRoles(args); err != nil {
		log.Fatalf("Fail to grant user 'user1' role 'writable' due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) revokeUserRoles() error {
	args := &api.RevokeUserRolesArgs{
		Username: "user1",
		Roles:    []string{"writable"},
	}
	if err := m.client.RevokeUserRoles(args); err != nil {
		log.Fatalf("Fail to grant user 'user1' role 'writable' due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) grantUserPrivileges() error {

	args := &api.GrantUserPrivilegesArgs{
		Username: "user1",
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "*",
				Table:      "*",
				Privileges: []string{"QUERY", "SELECT", "SEARCH"},
			},
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_CONTROL", "TABLE_READWRITE"},
			},
		},
	}
	if err := m.client.GrantUserPrivileges(args); err != nil {
		log.Fatalf("Fail to GrantUserPrivileges due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) revokeUserPrivileges() error {

	args := &api.RevokeUserPrivilegesArgs{
		Username: "user1",
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_READWRITE"},
			},
		},
	}
	if err := m.client.RevokeUserPrivileges(args); err != nil {
		log.Fatalf("Fail to RevokeUserPrivileges due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) showUserPrivileges() error {

	ret, err := m.client.ShowUserPrivileges("user1")
	if err != nil {
		log.Fatalf("Fail to ShowUserPrivileges due to error: %v", err)
		return err
	}
	fmt.Println("Privileges: ", ret)

	return nil
}

func (m *MochowTest) changeUserPassword() error {
	if err := m.client.ChangeUserPassword("user1", "mochow@456"); err != nil {
		log.Fatalf("Fail to change user 'user1' password due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) SelectUser() error {
	args := &api.SelectUserArgs{
		PrivilegeTuples: []api.PrivilegeTuple{
			{
				Database:   "book",
				Table:      "book_segments",
				Privileges: []string{"TABLE_CONTROL"},
			},
		},
	}
	ret, err := m.client.SelectUser(args)
	if err != nil {
		log.Fatalf("Fail to SelectUser due to error: %v", err)
		return err
	}
	fmt.Println("User: ", ret)
	return nil
}

func (m *MochowTest) dropUser() error {
	if err := m.client.DropUser("user1"); err != nil {
		log.Fatalf("Fail to create user due to error: %v", err)
		return err
	}
	return nil
}

func (m *MochowTest) testBinaryVector() error {
	// 1. 建表
	tableName := "binary_vector_table"
	fieldDim := 16 // 16 bits = 2 bytes
	fields := []api.FieldSchema{
		{
			FieldName:    "id",
			FieldType:    api.FieldTypeString,
			PrimaryKey:   true,
			PartitionKey: true,
			NotNull:      true,
		},
		{
			FieldName: "vector",
			FieldType: api.FieldTypeBinaryVector,
			NotNull:   true,
			Dimension: uint32(fieldDim),
		},
	}
	createTableArgs := &api.CreateTableArgs{
		Database:    m.database,
		Table:       tableName,
		Description: "test binary vector",
		Replication: 3,
		Partition: &api.PartitionParams{
			PartitionType: api.HASH,
			PartitionNum:  3,
		},
		Schema: &api.TableSchema{
			Fields: fields,
		},
	}

	// 检查表是否存在
	tableExisted, err := m.client.HasTable(m.database, tableName)
	if err != nil {
		log.Fatalf("Fail to check table due to error: %v", err)
		return err
	}
	if tableExisted {
		if err := m.client.DropTable(m.database, tableName); err != nil {
			log.Printf("Fail to drop table due to error: %v", err)
			return err
		}
		// wait until drop finished
		for {
			time.Sleep(time.Second * 5)
			if _, err := m.client.DescTable(m.database, tableName); err != nil {
				if realErr, ok := err.(*client.BceServiceError); ok {
					if realErr.StatusCode == 404 || realErr.Code == int(api.TableNotExist) {
						log.Println("Table already dropped")
						break
					}
				}
			}
		}
	}

	if err := m.client.CreateTable(createTableArgs); err != nil {
		log.Fatalf("Fail to create binary vector table: %v", err)
		return err
	}
	// 等待表 ready
	for {
		time.Sleep(time.Second * 5)
		desc, err := m.client.DescTable(m.database, tableName)
		if err == nil && desc.Table.State == api.TableStateNormal {
			log.Println("Binary vector table create finished")
			break
		}
	}

	// 2. 写入200条数据
	rows := make([]api.Row, 200)
	for i := 0; i < 200; i++ {
		vector := make([]byte, 2)
		if i < 16 {
			// 设置第i位为1
			byteIdx := i / 8
			bitIdx := i % 8
			vector[byteIdx] = 1 << bitIdx
		} else {
			// 后续数据使用其他模式，如递增值
			vector[0] = byte(i % 256)
			vector[1] = byte(i / 256)
		}

		rows[i] = api.Row{
			Fields: map[string]interface{}{
				"id":     fmt.Sprintf("b%d", i),
				"vector": api.BinaryVector(vector),
			},
		}
	}

	upsertArgs := &api.UpsertRowArg{
		Database: m.database,
		Table:    tableName,
		Rows:     rows,
	}
	result, err := m.client.UpsertRow(upsertArgs)
	if err != nil {
		log.Fatalf("Fail to upsert binary vector: %v", err)
		return err
	}
	log.Printf("Upserted %d binary vector rows", result.AffectedCount)

	// 3. 检索
	queryVec := api.BinaryVector([]byte{0x01, 0x00}) // 16维，第0位为1
	searchArgs := &api.VectorSearchArgs{
		Database: m.database,
		Table:    tableName,
		Request: api.VectorTopkSearchRequest{}.
			New("vector", queryVec, 10),
	}
	searchResult, err := m.client.VectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search binary vector: %v", err)
		return err
	}
	log.Printf("Binary vector search result: found %d rows", len(searchResult.Rows.Rows))
	m.client.DropTable(m.database, tableName)
	return nil
}

func (m *MochowTest) testSparseVector() error {
	// 1. 建表
	tableName := "sparse_vector_table"
	fields := []api.FieldSchema{
		{
			FieldName:    "id",
			FieldType:    api.FieldTypeString,
			PrimaryKey:   true,
			PartitionKey: true,
			NotNull:      true,
		},
		{
			FieldName: "vector",
			FieldType: api.FieldTypeSparseVector,
			NotNull:   true,
			Dimension: 1000, // 假设最大1000维
		},
	}
	indexes := []api.IndexSchema{
		{
			IndexName:  "sparse_vector_idx",
			Field:      "vector",
			IndexType:  api.SPARSE,
			MetricType: api.IP, // 稀疏向量常用L2
		},
	}
	createTableArgs := &api.CreateTableArgs{
		Database:    m.database,
		Table:       tableName,
		Description: "test sparse vector",
		Replication: 3,
		Partition: &api.PartitionParams{
			PartitionType: api.HASH,
			PartitionNum:  1,
		},
		Schema: &api.TableSchema{
			Fields:  fields,
			Indexes: indexes,
		},
	}

	// 检查表是否存在
	tableExisted, err := m.client.HasTable(m.database, tableName)
	if err != nil {
		log.Fatalf("Fail to check table due to error: %v", err)
		return err
	}
	if tableExisted {
		if err := m.client.DropTable(m.database, tableName); err != nil {
			log.Printf("Fail to drop table due to error: %v", err)
			return err
		}
		// wait until drop finished
		for {
			time.Sleep(time.Second * 5)
			if _, err := m.client.DescTable(m.database, tableName); err != nil {
				if realErr, ok := err.(*client.BceServiceError); ok {
					if realErr.StatusCode == 404 || realErr.Code == int(api.TableNotExist) {
						log.Println("Table already dropped")
						break
					}
				}
			}
		}
	}

	if err := m.client.CreateTable(createTableArgs); err != nil {
		log.Fatalf("Fail to create sparse vector table: %v", err)
		return err
	}
	for {
		time.Sleep(time.Second * 5)
		desc, err := m.client.DescTable(m.database, tableName)
		if err == nil && desc.Table.State == api.TableStateNormal {
			log.Println("Sparse vector table create finished")
			break
		}
	}

	// 2. 写入200条数据
	rows := make([]api.Row, 200)
	for i := 0; i < 200; i++ {
		nonZeroCount := 3 + (i % 3) // 3到5个非零元素
		sparseVec := make(api.SparseFloatVector)
		for j := 0; j < nonZeroCount; j++ {
			index := fmt.Sprintf("%d", (i*10+j*100)%1000)
			value := float32(0.1 + float64(i+j)/1000.0)
			sparseVec[index] = value
		}

		rows[i] = api.Row{
			Fields: map[string]interface{}{
				"id":     fmt.Sprintf("s%d", i),
				"vector": sparseVec,
			},
		}
	}

	upsertArgs := &api.UpsertRowArg{
		Database: m.database,
		Table:    tableName,
		Rows:     rows,
	}
	result, err := m.client.UpsertRow(upsertArgs)
	if err != nil {
		log.Fatalf("Fail to upsert sparse vector: %v", err)
		return err
	}
	log.Printf("Upserted %d sparse vector rows", result.AffectedCount)

	// 3. 检索
	queryVec := api.SparseFloatVector{"0": 0.1, "100": 0.2, "200": 0.3}
	searchArgs := &api.VectorSearchArgs{
		Database: m.database,
		Table:    tableName,
		Request: api.VectorTopkSearchRequest{}.
			New("vector", queryVec, 10),
	}
	searchResult, err := m.client.VectorSearch(searchArgs)
	if err != nil {
		log.Fatalf("Fail to search sparse vector: %v", err)
		return err
	}
	log.Printf("Sparse vector search result: found %d rows", len(searchResult.Rows.Rows))
	m.client.DropTable(m.database, tableName)
	return nil
}

func main() {
	// Init client
	clientConfig := &mochow.ClientConfiguration{
		Account:  "root",
		APIKey:   "*********",
		Endpoint: "http://*.*.*.*:*", // example:http://127.0.0.1:8511
	}
	mochowTest, err := NewMochowTest("book", "book_segments", clientConfig)
	if err != nil {
		log.Fatalf("Fail to init mochow example, err:%v", err)
	}
	log.Println("Init MochowTest success.")

	if err := mochowTest.clearEnv(); err != nil {
		fmt.Printf("Fail to clear env, err:%v", err)
		return
	}
	log.Println("Clear Env success.")

	// Create database and table
	if err := mochowTest.createDatabaseAndTable(); err != nil {
		log.Printf("Fail to create database and table, err:%v", err)
		return
	}
	log.Println("Create database and table success.")

	// Upsert data
	if err := mochowTest.upsertData(); err != nil {
		log.Printf("Fail to upsert row, err:%v", err)
		return
	}
	log.Println("Upsert row success.")

	// Query data
	if err := mochowTest.queryData(); err != nil {
		log.Printf("Fail to query row, err:%v", err)
		return
	}
	log.Println("Query row success.")

	// Batch query data
	if err := mochowTest.batchQueryData(); err != nil {
		log.Printf("Fail to query row, err:%v", err)
		return
	}
	log.Println("Batch query row success.")

	// Select data
	if err := mochowTest.selectData(); err != nil {
		log.Printf("Fail to select row, err:%v", err)
		return
	}
	log.Println("Select row success.")

	// Update data
	if err := mochowTest.updateData(); err != nil {
		log.Printf("Fail to update row, err:%v", err)
		return
	}
	log.Println("Update row success.")

	// Search data
	if err := mochowTest.topkSearch(); err != nil {
		log.Printf("Fail to topk search, err:%v", err)
		return
	}
	if err := mochowTest.rangeSearch(); err != nil {
		log.Printf("Fail to range search, err:%v", err)
		return
	}
	if err := mochowTest.batchSearch(); err != nil {
		log.Printf("Fail to batch search, err:%v", err)
		return
	}
	if err := mochowTest.searchIterator(); err != nil {
		log.Printf("Fail to use search iterator, err:%v", err)
		return
	}
	if err := mochowTest.bm25Search(); err != nil {
		log.Printf("Fail to bm25 search, err:%v", err)
		return
	}
	if err := mochowTest.hybridSearch(); err != nil {
		log.Printf("Fail to hybrid search, err:%v", err)
		return
	}
	if err := mochowTest.multiVectorSearch(); err != nil {
		log.Printf("Fail to multi vector search, err:%v", err)
		return
	}
	log.Println("Search row success.")

	// Delete data with pk
	if err := mochowTest.deleteDataWithPK(); err != nil {
		log.Printf("Fail to delete row with pk, err:%v", err)
		return
	}
	log.Println("Delete row with pk success.")

	// Delete data with filter
	if err := mochowTest.deleteDataWithFilter(); err != nil {
		log.Printf("Fail to delete row with filter, err:%v", err)
		return
	}
	log.Println("Delete row with filter success.")

	// drop and recreate vindex
	if err := mochowTest.dropAndCreateVIndex(); err != nil {
		log.Printf("Fail to drop and recreate vindex, err:%v", err)
		return
	}
	log.Println("Drop and recreate vindex success.")

	/////////////////// the following code is the example of rbac
	// role
	{
		if err := mochowTest.createRole(); err != nil {
			log.Printf("Fail to create role 'writable', err:%v", err)
			return
		}
		log.Println("Create role 'writable' success.")

		if err := mochowTest.grantRolePrivileges(); err != nil {
			log.Printf("Fail to grant role 'writable' privileges, err:%v", err)
			return
		}
		log.Println("Grant role 'writable' privileges success.")

		if err := mochowTest.showRolePrivileges(); err != nil {
			log.Printf("Fail to show role 'writable' privileges, err:%v", err)
			return
		}
		log.Println("Show role 'writable' privileges success.")

		if err := mochowTest.revokeRolePrivileges(); err != nil {
			log.Printf("Fail to revoke role 'writable' privileges, err:%v", err)
			return
		}
		log.Println("Revoke role 'writable' privileges success.")

		if err := mochowTest.showRolePrivileges(); err != nil {
			log.Printf("Fail to show role 'writable' privileges, err:%v", err)
			return
		}
		log.Println("Show role 'writable' privileges success.")

		if err := mochowTest.SelectRole(); err != nil {
			log.Printf("Fail to select role, err:%v", err)
			return
		}
		log.Println("Select rolw success.")
	}
	// user
	{
		if err := mochowTest.createUser(); err != nil {
			log.Printf("Fail to create user 'user1', err:%v", err)
			return
		}
		log.Println("Create user 'user1' success.")

		if err := mochowTest.grantUserRoles(); err != nil {
			log.Printf("Fail to grant user 'user1' role 'writable', err:%v", err)
			return
		}
		log.Println("Grant user 'user1' role 'writable' success.")

		if err := mochowTest.revokeUserRoles(); err != nil {
			log.Printf("Fail to revoke user 'user1' role 'writable', err:%v", err)
			return
		}
		log.Println("Revoke user 'user1' role 'writable' success.")

		if err := mochowTest.grantUserPrivileges(); err != nil {
			log.Printf("Fail to grant user 'user1' privileges, err:%v", err)
			return
		}
		log.Println("Show user 'user1' privileges success.")

		if err := mochowTest.revokeUserPrivileges(); err != nil {
			log.Printf("Fail to revoke user 'user1' privileges, err:%v", err)
			return
		}
		log.Println("Revoke user 'user1' privileges success.")

		if err := mochowTest.showUserPrivileges(); err != nil {
			log.Printf("Fail to show user 'user1' privileges, err:%v", err)
			return
		}
		log.Println("Show user 'user1' privileges success.")

		if err := mochowTest.SelectUser(); err != nil {
			log.Printf("Fail to select user, err:%v", err)
			return
		}
		log.Println("Select user success.")

		if err := mochowTest.changeUserPassword(); err != nil {
			log.Printf("Fail to change user 'user1' password, err:%v", err)
			return
		}
		log.Println("Change user 'user1' password success.")
	}

	if err := mochowTest.dropRole(); err != nil {
		log.Printf("Fail to drop role 'writable', err:%v", err)
		return
	}
	log.Println("Drop role 'writable' success.")

	if err := mochowTest.dropUser(); err != nil {
		log.Printf("Fail to drop user 'user1', err:%v", err)
		return
	}
	log.Println("Drop user 'user1' success.")

	// clear env
	if err := mochowTest.clearEnv(); err != nil {
		fmt.Printf("Fail to clear env, err:%v", err)
		return
	}
	log.Println("Clear Env success")

	// test binary vector
	if err := mochowTest.testBinaryVector(); err != nil {
		log.Printf("Fail to test binary vector, err:%v", err)
		return
	}
	log.Println("Test binary vector success.")

	// test sparse vector
	if err := mochowTest.testSparseVector(); err != nil {
		log.Printf("Fail to test sparse vector, err:%v", err)
		return
	}
	log.Println("Test sparse vector success.")
}
