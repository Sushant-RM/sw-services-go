package dto

// ============================================================
// DIGIT CORE REQUEST/RESPONSE
// ============================================================

type RequestInfo struct {
	APIID         string    `json:"apiId"`
	Ver           string    `json:"ver"`
	Ts            int64     `json:"ts"`
	Action        string    `json:"action"`
	Did           string    `json:"did,omitempty"`
	Key           string    `json:"key,omitempty"`
	MsgID         string    `json:"msgId"`
	AuthToken     string    `json:"authToken"`
	Correlationid string    `json:"correlationId,omitempty"`
	UserInfo      *UserInfo `json:"userInfo,omitempty"`
}

type UserInfo struct {
	ID           int64  `json:"id"`
	UUID         string `json:"uuid"`
	UserName     string `json:"userName"`
	Name         string `json:"name"`
	MobileNumber string `json:"mobileNumber"`
	EmailId      string `json:"emailId,omitempty"`
	Locale       string `json:"locale,omitempty"`
	Type         string `json:"type"`
	Roles        []Role `json:"roles"`
	Active       bool   `json:"active"`
	TenantID     string `json:"tenantId"`
}

type Role struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	TenantID string `json:"tenantId"`
}

type ResponseInfo struct {
	APIID    string `json:"apiId"`
	Ver      string `json:"ver"`
	Ts       int64  `json:"ts"`
	ResMsgID string `json:"resMsgId"`
	MsgID    string `json:"msgId"`
	Status   string `json:"status"` // "successful" or "failed"
}

func NewResponseInfo(reqInfo *RequestInfo, status string) ResponseInfo {
	return ResponseInfo{
		APIID:  reqInfo.APIID,
		Ver:    reqInfo.Ver,
		Ts:     reqInfo.Ts,
		MsgID:  reqInfo.MsgID,
		Status: status,
	}
}

// ============================================================
// IDGEN STRUCTS
// ============================================================

type IDGenRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
	IDRequests  []IDRequest `json:"idRequests"`
}

type IDRequest struct {
	IDName   string `json:"idName"`
	Format   string `json:"format"`
	TenantID string `json:"tenantId"`
	Count    int    `json:"count"`
}

type IDGenResponse struct {
	ResponseInfo ResponseInfo `json:"ResponseInfo"`
	IDResponses  []IDResponse `json:"idResponses"`
}

type IDResponse struct {
	IDName string `json:"idName"`
	Format string `json:"format"`
	ID     string `json:"id"`
}

// ============================================================
// WORKFLOW STRUCTS
// ============================================================

type WorkflowRequest struct {
	RequestInfo      RequestInfo       `json:"RequestInfo"`
	ProcessInstances []ProcessInstance `json:"ProcessInstances"`
}

type ProcessInstance struct {
	BusinessService string     `json:"businessService"`
	BusinessID      string     `json:"businessId"`
	TenantID        string     `json:"tenantId"`
	Action          string     `json:"action"`
	Comment         string     `json:"comment,omitempty"`
	Assignes        []Assignee `json:"assignes,omitempty"`
	Documents       []Document `json:"documents,omitempty"`
	ModuleName      string     `json:"moduleName"`
}

type Assignee struct {
	UUID string `json:"uuid"`
}

type Document struct {
	ID           string        `json:"id,omitempty"`
	DocumentType string        `json:"documentType"`
	FileStoreId  string        `json:"fileStoreId"`
	DocumentUID  string        `json:"documentUid,omitempty"`
	AuditDetails *AuditDetails `json:"auditDetails,omitempty"`
	Status       string        `json:"status,omitempty"`
}

type WorkflowResponse struct {
	ResponseInfo     ResponseInfo      `json:"ResponseInfo"`
	ProcessInstances []ProcessInstance `json:"ProcessInstances"`
}

// ============================================================
// SEWERAGE CONNECTION
// ============================================================

type PlumberInfo struct {
	ID                    string        `json:"id"`
	Name                  string        `json:"name"`
	LicenseNo             string        `json:"licenseNo"`
	MobileNumber          string        `json:"mobileNumber"`
	Gender                string        `json:"gender"`
	FatherOrHusbandName   string        `json:"fatherOrHusbandName"`
	CorrespondenceAddress string        `json:"correspondenceAddress"`
	Relationship          string        `json:"relationship"`
	AuditDetails          *AuditDetails `json:"auditDetails"`
}

type RoadCuttingInfo struct {
	ID              string        `json:"id"`
	Status          string        `json:"status"`
	RoadType        string        `json:"roadType"`
	RoadCuttingArea int           `json:"roadCuttingArea"`
	AuditDetails    *AuditDetails `json:"auditDetails"`
}

type SewerageConnection struct {
	ID                         string             `json:"id,omitempty"`
	TenantID                   string             `json:"tenantId"`
	ConnectionNo               string             `json:"connectionNo"`
	OldConnectionNo            string             `json:"oldConnectionNo"`
	ApplicationNumber          string             `json:"applicationNumber,omitempty"` // for Go repository compatibility
	ApplicationNo              string             `json:"applicationNo"`     // DIGIT exact contract
	ApplicationType            string             `json:"applicationType"`             // "NEW_CONNECTION"
	ApplicationStatus          string             `json:"applicationStatus"` // Go service DB representation
	Status                     string             `json:"status"`            // DIGIT exact representation
	ConnectionType             string             `json:"connectionType"`              // "Permanent"
	NoOfToilets                int                `json:"noOfToilets,omitempty"`
	NoOfWaterClosets           int                `json:"noOfWaterClosets,omitempty"`
	PropertyID                 string             `json:"propertyId,omitempty"`
	PropertyUsageType          string             `json:"propertyUsageType"`
	ConnectionHolders          []ConnectionHolder `json:"connectionHolders"`
	PlumberInfo                []PlumberInfo      `json:"plumberInfo"`
	RoadCuttingInfo            []RoadCuttingInfo  `json:"roadCuttingInfo"`
	Documents                  []Document         `json:"documents"`
	ProcessInstance            ProcessInstance    `json:"processInstance"`
	AdditionalDetails          interface{}        `json:"additionalDetails"`
	AuditDetails               *AuditDetails      `json:"auditDetails,omitempty"`
	Channel                    string             `json:"channel"`
	RoadType                   string             `json:"roadType"`
	RoadCuttingArea            int                `json:"roadCuttingArea,omitempty"`
	OldApplication             bool               `json:"oldApplication"`
	IsDisconnectionTemporary   bool               `json:"isDisconnectionTemporary"`
	DisconnectionReason        string             `json:"disconnectionReason"`
	DateEffectiveFrom          int64              `json:"dateEffectiveFrom"`
	ConnectionExecutionDate    int64              `json:"connectionExecutionDate"`
	ProposedWaterClosets       int                `json:"proposedWaterClosets"`
	ProposedToilets            int                `json:"proposedToilets"`
	DisconnectionExecutionDate int64              `json:"disconnectionExecutionDate"`
}

type ConnectionHolder struct {
	ID                    string  `json:"id,omitempty"`
	TenantID              string  `json:"tenantId"`
	OwnerShipCategory     string  `json:"ownerShipCategory,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Gender                string  `json:"gender,omitempty"`
	FatherOrHusbandName   string  `json:"fatherOrHusbandName,omitempty"`
	Relationship          string  `json:"relationship"`
	MobileNumber          string  `json:"mobileNumber,omitempty"`
	IsSameAsPropertyOwner bool    `json:"isSameAsPropertyOwner,omitempty"`
	UUID                  string  `json:"uuid"`
	OwnerType             string  `json:"ownerType,omitempty"`
	Status                string  `json:"status"`
	IsPrimaryOwner        bool    `json:"isPrimaryOwner"`
	OwnerShipPercentage   float64 `json:"ownerShipPercentage"`
}

type AuditDetails struct {
	CreatedBy        string `json:"createdBy"`
	LastModifiedBy   string `json:"lastModifiedBy"`
	CreatedTime      int64  `json:"createdTime"`
	LastModifiedTime int64  `json:"lastModifiedTime"`
}

type CreateSewerageConnectionRequest struct {
	RequestInfo                RequestInfo        `json:"RequestInfo"`
	SewerageConnection         SewerageConnection `json:"SewerageConnection"`
	IsCreateCall               bool               `json:"isCreateCall,omitempty"`
	IsOldDataEncryptionRequest bool               `json:"isOldDataEncryptionRequest,omitempty"`
	DisconnectRequest          bool               `json:"disconnectRequest,omitempty"`
}

type SewerageConnectionRequest struct {
	RequestInfo                RequestInfo        `json:"RequestInfo"`
	SewerageConnection         SewerageConnection `json:"SewerageConnection"`
	IsCreateCall               bool               `json:"isCreateCall,omitempty"`
	IsOldDataEncryptionRequest bool               `json:"isOldDataEncryptionRequest,omitempty"`
	DisconnectRequest          bool               `json:"disconnectRequest,omitempty"`
}

type SewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
	TotalCount          int                  `json:"totalCount,omitempty"`
}

type CreateSewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
}

type SearchSewerageConnectionResponse struct {
	ResponseInfo        ResponseInfo         `json:"ResponseInfo"`
	SewerageConnections []SewerageConnection `json:"SewerageConnections"`
	TotalCount          int                  `json:"TotalCount"`
}

type SearchSewerageConnectionRequest struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}

type RequestInfoWrapper struct {
	RequestInfo RequestInfo `json:"RequestInfo"`
}
