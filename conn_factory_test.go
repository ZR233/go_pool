package go_pool

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestGetHostsFuncFromConsul(t *testing.T) {

	client, err := api.NewClient(api.DefaultConfig()) //非默认情况下需要设置实际的参数
	if err != nil {
		t.Error(err)
		return
	}
	f := GetHostsFuncFromConsul(client, "digger/packet_check")
	hosts, err := f()
	println(hosts)

}
