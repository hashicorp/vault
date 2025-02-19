# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1

storage "raft" {
  path    = "/storage/path/raft"
  node_id = "raft1"

  retry_join {
    "leader_api_addr" = "http://127.0.0.1:8200"
  }
  retry_join {
    "leader_api_addr" = "http://[2001:db8:0:0:0:0:2:1]:8200"
  }
  retry_join {
    "auto_join" = "provider=mdns service=consul domain=2001:db8:0:0:0:0:2:1"
  }
  retry_join {
    "auto_join" = "provider=os tag_key=consul tag_value=server username=foo password=bar auth_url=https://[2001:db8:0:0:0:0:2:1]/auth"
  }
  retry_join {
    "auto_join" = "provider=triton account=testaccount url=https://[2001:db8:0:0:0:0:2:1] key_id=1234 tag_key=consul-role tag_value=server"
  }
  retry_join {
    "auto_join" = "provider=packet auth_token=token project=uuid url=https://[2001:db8:0:0:0:0:2:1] address_type=public_v6"
  }
  retry_join {
    "auto_join" = "provider=vsphere category_name=consul-role tag_name=consul-server host=https://[2001:db8:0:0:0:0:2:1] user=foo password=bar insecure_ssl=false"
  }
}

listener "tcp" {
  address = "127.0.0.1:8200"
}

disable_mlock = true
