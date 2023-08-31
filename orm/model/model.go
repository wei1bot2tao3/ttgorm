package model

import (
	"reflect"
	"strings"
	"sync"
	"ttgorm/orm/internal/errs"
	"unicode"
)

const (
	togColumn = "column"
)

type Registry interface {
	Get(val any) (*Model, error)
	Registry(val any, opts ...Option) (*Model, error)
}

type Option func(m *Model) error

func WithTableName(teblename string) Option {
	return func(m *Model) error {
		m.TableName = teblename
		//if teblename==""{
		//	return errors.New("")
		//}
		return nil
	}
}

func WithColumnName(field string, column string) Option {
	return func(m *Model) error {
		fd, ok := m.FieldsMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = column
		return nil
	}
}

// Model 表示一个注册在全局的数据模型元数据，包含表名和对应的列名。
// Filedsmap是 字段名到字段到映射 ColumnMP是数据库的列名到字段的映射
type Model struct {
	TableName string            // 表名
	Fields    []*Field          // 提前计算好的列名和对应的字段
	FieldsMap map[string]*Field // 字段名到字段的映射
	ColumnMap map[string]*Field // 数据库列名到字段的映射
}

// Field 表示一个字段
type Field struct {
	// 列名
	ColName string

	//代表 字段类型
	Type reflect.Type

	//字段名
	GoName string
	// 偏移量
	Offset uintptr
}

// UnderscoreName 驼峰转字符串命名
func UnderscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}

// registry 元数据的注册中心
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	// reflect.Type 唯一的
	//models map[reflect.Type]*Model
	// models 是存放 表名和 列名的
	models sync.Map
}

// NewRegistry 获取一个registry实例
func NewRegistry() Registry {
	return &registry{
		//models: make(map[reflect.Type]*Model, 64),
	}
}

// Get 获取 加注册  这个val 是对应的表 在go的结构体
func (r *registry) Get(val any) (*Model, error) {
	// 获取 它对的类型
	typ := reflect.TypeOf(val)
	// 从注册表中加载对应的数据模型  没有的和 在去创建 下次就有了
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Registry(val)
	if err != nil {
		return nil, err
	}

	return m.(*Model), err
}

// Registry 通过传入指向结构体的指针 结构体--》元数据的过程
func (r *registry) Registry(entity any, opts ...Option) (*Model, error) {
	// 先获取 结构体的类型
	typ := reflect.TypeOf(entity)
	// 判断是否传入指向结构体的指针
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	// 指针的话 先elem 取对应字段
	elemType := typ.Elem()
	//获取长度啊
	numFiled := elemType.NumField()
	// 创建 go结构体到元数据的映射
	fieldMap := make(map[string]*Field, numFiled)
	columnMap := make(map[string]*Field, numFiled)
	fields := make([]*Field, 0, numFiled)
	for i := 0; i < numFiled; i++ {
		filedType := elemType.Field(i)
		pair, err := r.parseTag(filedType.Tag)
		if err != nil {
			return nil, err
		}
		columnName := pair[togColumn]
		if columnName == "" {
			// 用户没有设置标签 就取这个结构体的 字段名
			columnName = UnderscoreName(filedType.Name)
		}
		// 构建一下 go结构体下 每一个字段 的信息  把它存到file map里
		fdMeta := &Field{
			ColName: columnName,
			Type:    filedType.Type,
			GoName:  filedType.Name,
			Offset:  filedType.Offset,
		}
		// 根据字段名 来找 字段信息
		fieldMap[filedType.Name] = fdMeta
		// 根据数据库重的列名找
		columnMap[columnName] = fdMeta
		fields = append(fields, fdMeta)

	}

	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}
	// 用户没有设置
	if tableName == "" {
		tableName = UnderscoreName(elemType.Name())
	}
	res := &Model{
		TableName: tableName,
		ColumnMap: columnMap,
		FieldsMap: fieldMap,
		Fields:    fields,
	}
	for _, opt := range opts {

		err := opt(res)
		if err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, res)

	return res, nil
}

// parseTag 返回一个标签
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		value := segs[1]
		res[key] = value
	}

	return res, nil
}

type User struct {
	ID uint64 `orm:"column=id,xxx=bbb"`
}

type TableName interface {
	TableName() string
}
