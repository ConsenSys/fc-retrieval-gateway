package messages

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// admin.go contains all messages originting from the gateway related to the admin API

// AdminGetReputationChallenge is the request from an admin client to a gateway to discover a client's reputation
type AdminGetReputationChallenge struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	ClientID        string `json:"clientid"`
}

// AdminGetReputationResponse is the response to AdminGetReputationChallenge
type AdminGetReputationResponse struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	ClientID        string `json:"clientid"`
	Reputation      int64  `json:"reputation"`
	Exists          bool   `json:"exists"`
}

// AdminSetReputationChallenge is the request from an admin client to a gateway to set a client's reputation
type AdminSetReputationChallenge struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	ClientID        string `json:"clientid"`
	Reputation      int64  `json:"reputation"`
}

// AdminSetReputationResponse is the response to AdminSetReputationChallenge
type AdminSetReputationResponse struct {
	MessageType     int32  `json:"message_type"`
	ProtocolVersion int32  `json:"protocol_version"`
	ClientID        string `json:"clientid"`
	Reputation      int64  `json:"reputation"`
	Exists          bool   `json:"exists"`
}
