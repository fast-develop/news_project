package handler

import (
    "sync"
	"context"
	//"errors"
	//"sort"
	//"strings"

	//"strconv"

	//"encoding/json"
	//"errors"
	//"fmt"
	log "github.com/golang/glog"

	"engine/proto/engine"
)

func HandleRequest(ctx context.Context, in *engine.EngineRequest) (*engine.EngineResponse, error) {
    //fmt.Println("HandleRequest")
	countScope := countstat.GetCountGlobal().NewCountScope("HandleRequest")
	countScope.SetErr()
	defer countScope.End()

	rsp := &engine.EngineResponse{}

    userId := in.GetUserId()
    docIdList := in.GetDocIdList()
    fromList := in.GetFromList()
    var finalMap sync.Map
    docMap := make(map[int32]int32)
    for i, docid := range docIdList {
        docMap[docid] = fromList[i]
    }
    
    log.Infof("handle request userid:%s", userId)

    //fmt.Printf("%v",finalMap)

    errId := int32(0)
	rsp.ErrorId = &errId
	//fmt.Println(rsp)

	return rsp, nil
}
