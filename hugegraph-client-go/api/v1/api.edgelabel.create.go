/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements. See the NOTICE file distributed with this
 * work for additional information regarding copyright ownership. The ASF
 * licenses this file to You under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package v1

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strings"

    "github.com/apache/incubator-hugegraph-toolchain/hugegraph-client-go/api"
)

// ----- API Definition -------------------------------------------------------
// Create an EdgeLabel
//
// See full documentation at https://hugegraph.apache.org/docs/clients/restful-api/edgelabel/#141-create-an-edgelabel
func newEdgelabelCreateFunc(t api.Transport) EdgelabelCreate {
    return func(o ...func(*EdgelabelCreateRequest)) (*EdgelabelCreateResponse, error) {
        var r = EdgelabelCreateRequest{}
        for _, f := range o {
            f(&r)
        }
        return r.Do(r.ctx, t)
    }
}

type EdgelabelCreate func(o ...func(*EdgelabelCreateRequest)) (*EdgelabelCreateResponse, error)

type EdgelabelCreateRequest struct {
    Body    io.Reader
    ctx     context.Context
    reqData EdgelabelCreateRequestData
}

type EdgelabelCreateRequestData struct {
    Name             string   `json:"name"`
    SourceLabel      string   `json:"source_label"`
    TargetLabel      string   `json:"target_label"`
    Frequency        string   `json:"frequency"`
    Properties       []string `json:"properties"`
    SortKeys         []string `json:"sort_keys"`
    NullableKeys     []string `json:"nullable_keys"`
    EnableLabelIndex bool     `json:"enable_label_index"`
}

type EdgelabelCreateResponse struct {
    StatusCode int                         `json:"-"`
    Header     http.Header                 `json:"-"`
    Body       io.ReadCloser               `json:"-"`
    RespData   EdgelabelCreateResponseData `json:"-"`
}

type EdgelabelCreateResponseData struct {
    ID               int      `json:"id"`
    SortKeys         []string `json:"sort_keys"`
    SourceLabel      string   `json:"source_label"`
    Name             string   `json:"name"`
    IndexNames       []string `json:"index_names"`
    Properties       []string `json:"properties"`
    TargetLabel      string   `json:"target_label"`
    Frequency        string   `json:"frequency"`
    NullableKeys     []string `json:"nullable_keys"`
    EnableLabelIndex bool     `json:"enable_label_index"`
    UserData         struct {
    } `json:"user_data"`
}

func (r EdgelabelCreateRequest) Do(ctx context.Context, transport api.Transport) (*EdgelabelCreateResponse, error) {

    if len(r.reqData.Name) <= 0 {
        return nil, errors.New("create edgeLabel must set name")
    }
    if len(r.reqData.SourceLabel) <= 0 {
        return nil, errors.New("create edgeLabel must set source_label")
    }
    if len(r.reqData.TargetLabel) <= 0 {
        return nil, errors.New("create edgeLabel must set target_label")
    }
    if len(r.reqData.Properties) <= 0 {
        return nil, errors.New("create edgeLabel must set properties")
    }

    byteBody, err := json.Marshal(&r.reqData)
    if err != nil {
        return nil, err
    }
    reader := strings.NewReader(string(byteBody))

    req, err := api.NewRequest("POST", fmt.Sprintf("/graphs/%s/schema/vertexlabels", transport.GetConfig().Graph), nil, reader)
    if err != nil {
        return nil, err
    }
    if ctx != nil {
        req = req.WithContext(ctx)
    }

    res, err := transport.Perform(req)
    if err != nil {
        return nil, err
    }

    resp := &EdgelabelCreateResponse{}
    respData := EdgelabelCreateResponseData{}
    bytes, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(bytes, &respData)
    if err != nil {
        return nil, err
    }
    resp.StatusCode = res.StatusCode
    resp.Header = res.Header
    resp.Body = res.Body
    resp.RespData = respData
    return resp, nil
}

func (r EdgelabelCreate) WithReqData(reqData EdgelabelCreateRequestData) func(request *EdgelabelCreateRequest) {
    return func(r *EdgelabelCreateRequest) {
        r.reqData = reqData
    }
}
