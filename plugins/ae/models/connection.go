/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

type AeConnection struct {
	AppId    string `mapstructure:"appId" env:"AE_APP_ID" json:"appId"`
	Sign     string `mapstructure:"sign" env:"AE_SIGN" json:"sign"`
	NonceStr string `mapstructure:"nonceStr" env:"AE_NONCE_STR" json:"nonceStr"`
	Endpoint string `mapstructure:"endpoint" env:"AE_ENDPOINT" json:"endpoint"`
}

// This object conforms to what the frontend currently expects.
type AeResponse struct {
	AeConnection
	Name string `json:"name"`
	ID   int    `json:"id"`
}
