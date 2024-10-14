package aliyun

type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type DomainStatus struct {
	DomainName string `json:"DomainName"`
	RecordId   string `json:"RecordId"`
	Value      string `json:"Value"`
}

type DomainStatusRes struct {
	TotalCount    uint `json:"TotalCount"`
	DomainRecords struct {
		Record []DomainStatus `json:"Record"`
	} `json:"DomainRecords"`
}
