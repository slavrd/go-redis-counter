# restrict access to READ redis password
path "kv/redispassword" {
  capabilities = ["read"]  
}