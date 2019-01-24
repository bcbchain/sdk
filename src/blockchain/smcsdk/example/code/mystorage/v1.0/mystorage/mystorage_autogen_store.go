package mystorage

//_setStoredData This is a method of MyStorage
func (ms *MyStorage) _setStoredData(v uint64) {
	ms.sdk.Helper().StateHelper().Set("/storedData", &v)
}

//_storedData This is a method of MyStorage
func (ms *MyStorage) _storedData() uint64 {

	return *ms.sdk.Helper().StateHelper().GetEx("/storedData", new(uint64)).(*uint64)
}

//_chkStoredData This is a method of MyStorage
func (ms *MyStorage) _chkStoredData() bool {
	return ms.sdk.Helper().StateHelper().Check("/storedData")
}
