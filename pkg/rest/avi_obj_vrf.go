/*
 * [2013] - [2018] Avi Networks Incorporated
 * All Rights Reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *   http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rest

import (
	"encoding/json"
	"errors"

	avicache "ako/pkg/cache"
	"ako/pkg/lib"
	"ako/pkg/nodes"

	"github.com/avinetworks/container-lib/utils"
	avimodels "github.com/avinetworks/sdk/go/models"
)

func (rest *RestOperations) AviVrfGet(key, uuid, name string) *avimodels.VrfContext {
	if rest.aviRestPoolClient == nil {
		utils.AviLog.Warnf("key: %s, msg: aviRestPoolClient not initialized\n", key)
		return nil
	}
	if len(rest.aviRestPoolClient.AviClient) < 1 {
		utils.AviLog.Warnf("key: %s, msg: client in aviRestPoolClient not initialized\n", key)
		return nil
	}
	client := rest.aviRestPoolClient.AviClient[0]
	uri := "/api/vrfcontext/" + uuid

	rawData, err := client.AviSession.GetRaw(uri)
	if err != nil {
		utils.AviLog.Warnf("Vrf Get uri %v returned err %v", uri, err)
		return nil
	}
	vrf := avimodels.VrfContext{}
	json.Unmarshal(rawData, &vrf)

	return &vrf
}

func (rest *RestOperations) AviVrfBuild(key string, vrfNode *nodes.AviVrfNode, uuid string) *utils.RestOp {
	vrfCacheObj := rest.getVrfCacheObj(vrfNode.Name)
	if vrfCacheObj == nil {
		return nil
	}
	vrf := rest.AviVrfGet(key, vrfCacheObj.Uuid, vrfCacheObj.Name)
	if vrf == nil {
		return nil
	}
	path := "/api/vrfcontext/" + vrfCacheObj.Uuid
	vrf.StaticRoutes = vrfNode.StaticRoutes
	restOp := utils.RestOp{Path: path, Method: utils.RestPut, Obj: vrf,
		Tenant: lib.GetTenant(), Model: "VrfContext", Version: utils.CtrlVersion}
	return &restOp
}

func (rest *RestOperations) getVrfCacheObj(vrfName string) *avicache.AviVrfCache {
	vrfCache, found := rest.cache.VrfCache.AviCacheGet(vrfName)
	if found {
		vrfCacheObj, ok := vrfCache.(*avicache.AviVrfCache)
		if !ok {
			utils.AviLog.Warnf("Vrf object for %s found. Cannot cast. Not doing anything\n", vrfName)
			return nil
		}
		return vrfCacheObj
	}
	utils.AviLog.Infof("vrf cache object NOT found for vrf name: %s", vrfName)
	return nil
}

func (rest *RestOperations) AviVrfCacheAdd(restOp *utils.RestOp, vrfKey avicache.NamespaceName, key string) error {
	if (restOp.Err != nil) || (restOp.Response == nil) {
		utils.AviLog.Warnf("key: %s, rest_op has err or no reponse for vrfcontext, err: %s, response: %s", key, restOp.Err, restOp.Response)
		return errors.New("Errored rest_op")
	}
	respElems, ok := RestRespArrToObjByType(restOp, "vrfcontext", key)
	if ok != nil || respElems == nil {
		utils.AviLog.Warnf("key: %s, msg: unable to find vrfcontext obj in resp %v", key, restOp.Response)
		return errors.New("vrfcontext not found")
	}
	vrfName := vrfKey.Name

	var checksum uint32
	var staticRoutes []*avimodels.StaticRoute
	rest.cache.VrfCache.AviCacheGet(vrfName)
	for _, resp := range respElems {
		name, ok := resp["name"].(string)
		if !ok {
			utils.AviLog.Warnf("key: %s, msg: wrong object type %T for name in vrf %s\n", key, resp["name"], vrfName)
			continue
		}
		uuid, ok := resp["uuid"].(string)
		if !ok {
			utils.AviLog.Warnf("key: %s, msg: wrong object type %T for uuid in vrf %s\n", key, resp["uuid"], vrfName)
			continue
		}
		if resp["static_routes"] == nil {
			staticRoutes = nil
		} else {
			staticRoutesIntf, ok := resp["static_routes"].([]interface{})
			if !ok {
				utils.AviLog.Warnf("key: %s, msg: wrong object type %T for staticroutes in staticroutes %s\n", key, resp["staticroutes"], vrfName)
				continue
			}
			staticRoutes = lib.StaticRoutesIntfToObj(staticRoutesIntf)
			if len(staticRoutes) == 0 {
				utils.AviLog.Debugf("key: %s, no static routes found for vrf %s\n", key, vrfName)
			}
		}
		checksum = lib.VrfChecksum(name, staticRoutes)
		vrfCacheObj := avicache.AviVrfCache{Name: name, Uuid: uuid, CloudConfigCksum: checksum}
		rest.cache.VrfCache.AviCacheAdd(vrfName, &vrfCacheObj)
	}
	return nil
}
