package v1

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
	Get(val any)(*Model,error)
	Registry(val any,opts ...ModelOption)(*Model,error)
}

type ModelOption func(m *Model)  error


// Model 注册在 全局的 数据模型
type Model struct {
	tableName string
	fields    map[string]*Field
}

// Field 表示一个字段
type Field struct {
	// 列名
	colName string
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
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
// reflect.Typ 唯一的
// models 是存放 表名和 列名的
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	//models map[reflect.Type]*Model
	models sync.Map
}

// 获取一个registry实例
func newRegistry() *registry {
	return &registry{
		//models: make(map[reflect.Type]*Model, 64),
	}
}

// Get
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Registry(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*Model), err
}

//func (r *registry) Get(val any) (*Model, error) {
//typ := reflect.TypeOf(val)
//r.lock.RLock()
//m, ok := r.models[typ]
//r.lock.RUnlock()
//if ok {
//	return m, nil
//}
//
//r.lock.Lock()
//defer r.lock.Unlock()
//m, ok = r.models[typ]
//if ok {
//	return m, nil
//}
//if !ok {
//	var err error
//	m, err = r.Registry(val)
//	if err != nil {
//		return nil, err
//	}
//
//	r.models[typ] = m
//}

//
//	return m, nil
//}

// Registry 通过传入指向结构体的指针 映射结构体的类型
func (r *registry) Registry(entity any,opts ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}

	typ = typ.Elem()
	//for typ.Kind() == reflect.Pointer {
	//	typ = typ.Elem()
	//}

	numFiled := typ.NumField()
	fieldMap := make(map[string]*Field, numFiled)
	for i := 0; i < numFiled; i++ {
		filedType := typ.Field(i)
		pair, err := r.parseTag(filedType.Tag)
		if err != nil {
			return nil, err
		}
		columnName := pair[togColumn]
		if columnName == "" {
			// 用户没有设置
			columnName = underscoreName(filedType.Name)
		}

		fieldMap[filedType.Name] = &Field{
			colName: columnName,
		}
	}

	var tablename string
	if tbl, ok := entity.(TableName); ok {
		tablename = tbl.TableName()
	}
	if tablename == "" {
		tablename = underscoreName(typ.Name())
	}
	res:=&Model{

		tableName: tablename,

		fields: fieldMap,
	}

	for _,opt:=range opts{
		opt(res)
	}

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


