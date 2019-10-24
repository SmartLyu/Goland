package Global

import "sync"

// 存储报警信息至内存
type ErrorMapType struct {
	Data map[string]int
	Lock sync.Mutex
}

// 存储Nat机器中子服务器信息至内存
type NatHostsMapType struct {
	Data map[HostsTable]bool
	Lock sync.Mutex
}

// 初始化创建map
func NewErrorMapType() *ErrorMapType{
	return &ErrorMapType{
		Data: make(map[string]int),
	}
}

// 初始化创建map
func NewNatHostsMapType() *NatHostsMapType{
	return &NatHostsMapType{
		Data: make(map[HostsTable]bool),
	}
}

// 是否存在
func (m *ErrorMapType) Exist(key string) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if m.Data[key] > 0 {
		return true
	} else {
		return false
	}
}

// 获取指定值
func (m *ErrorMapType) Get(key string) int {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	return m.Data[key]
}

// 获取所有值
func (m *ErrorMapType) Getall() (Keys []string, Values []int) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	for key, value := range m.Data {
		Keys = append(Keys, key)
		Values = append(Values, value)
	}
	return
}

// 值增加1
func (m *ErrorMapType) Add(key string, num int) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Data[key] = m.Data[key] + num
}

// 删除指定值
func (m *ErrorMapType) Delete(key string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	delete(m.Data, key)
}

// 功能同上
func (m *NatHostsMapType) Exist(ht HostsTable) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if _, isError := m.Data[ht]; isError {
		return true
	} else {
		return false
	}
}

func (m *NatHostsMapType) Get(ht HostsTable) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	return m.Data[ht]
}

func (m *NatHostsMapType) Delete(ht HostsTable) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	delete(m.Data, ht)
}

// 值状态更改
func (m *NatHostsMapType) Change(ht HostsTable, status bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Data[ht] = status
}

// 获取所有值
func (m *NatHostsMapType) Getall() (keys []HostsTable, values []bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	for key, value := range m.Data {
		keys = append(keys, key)
		values = append(values, value)
	}
	return
}
