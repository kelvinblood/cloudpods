package hw_fusion

const (
	BASIC_URI = `/service`
	LOGIN_URI = BASIC_URI + `/session`
	SITES_URI = BASIC_URI + `/sites`

	PREFIX_SITE_URI      = SITES_URI + `/<site_id>` // 可替换的URI
	VMS_URI              = PREFIX_SITE_URI + `/vms`
	VMS_STATISTICS_URI   = VMS_URI + `/statistics`
	HOSTS_URI            = PREFIX_SITE_URI + `/hosts`
	HOSTS_STATISTICS_URI = HOSTS_URI + `/statistics`
)
