package Global

import "sync"

// 存储报警信息至内存
type ErrorMapType struct {
	Data map[string]int
	Lock sync.Mutex
}

// 初始化创建map
func NewErrorMapType() *ErrorMapType {
	return &ErrorMapType{
		Data: make(map[string]int),
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

// 存储状态和hostname的信息
type StatusAndHostNameJson struct {
	Status   bool
	Hostname string
}

// 存储Nat机器中子服务器信息至内存
type NatHostsMapType struct {
	Data map[string]StatusAndHostNameJson
	Lock sync.Mutex
}

// 初始化创建map
func NewNatHostsMapType() *NatHostsMapType {
	return &NatHostsMapType{
		Data: make(map[string]StatusAndHostNameJson),
	}
}

// 功能同上
func (m *NatHostsMapType) Exist(ht HostsTable) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	if _, isError := m.Data[ht.IP]; isError {
		return true
	} else {
		return false
	}
}

func (m *NatHostsMapType) Get(ht HostsTable) StatusAndHostNameJson {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	return m.Data[ht.IP]
}

func (m *NatHostsMapType) Delete(ht HostsTable) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	delete(m.Data, ht.IP)
}

// 值状态更改
func (m *NatHostsMapType) Change(ht HostsTable, status bool) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	var tmpStatus = StatusAndHostNameJson{
		Status:   status,
		Hostname: ht.HostName,
	}
	m.Data[ht.IP] = tmpStatus
}
