package canoe

import (
	"bytes"
	"encoding/json"
	"fmt"
	cTypes "github.com/compose/canoe/types"
	eTypes "github.com/coreos/etcd/pkg/types"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

var peerEndpoint = "/peers"

func (rn *Node) peerAPI() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(peerEndpoint, rn.peerAddHandlerFunc()).Methods("POST")
	r.HandleFunc(peerEndpoint, rn.peerDeleteHandlerFunc()).Methods("DELETE")
	r.HandleFunc(peerEndpoint, rn.peerMembersHandlerFunc()).Methods("GET")

	return r
}

func (rn *Node) serveHTTP() error {
	router := rn.peerAPI()

	ln, err := newStoppableListener(fmt.Sprintf(":%d", rn.configPort), rn.stopc)
	if err != nil {
		panic(err)
	}

	err = (&http.Server{Handler: router}).Serve(ln)
	select {
	case <-rn.stopc:
		return nil
	default:
		return errors.Wrap(err, "Error serving HTTP API")
	}
}

func (rn *Node) serveRaft() error {
	ln, err := newStoppableListener(fmt.Sprintf(":%d", rn.raftPort), rn.stopc)
	if err != nil {
		return errors.Wrap(err, "Error creating a new stoppable listener")
	}

	err = (&http.Server{Handler: rn.transport.Handler()}).Serve(ln)

	select {
	case <-rn.stopc:
		return nil
	default:
		return errors.Wrap(err, "Error serving raft http server")
	}
}

func (rn *Node) peerMembersHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rn.handlePeerMembersRequest(w, req)
	}
}

func (rn *Node) handlePeerMembersRequest(w http.ResponseWriter, req *http.Request) {
	if !rn.initialized {
		rn.writeNodeNotReady(w)
	} else {
		membersResp := &cTypes.ConfigMembershipResponseData{
			cTypes.ConfigPeerData{
				RaftPort:          rn.raftPort,
				ConfigurationPort: rn.configPort,
				ID:                rn.id,
				RemotePeers:       rn.peerMap,
			},
		}

		rn.writeSuccess(w, membersResp)
	}
}

func (rn *Node) peerDeleteHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rn.handlePeerDeleteRequest(w, req)
	}
}

func (rn *Node) handlePeerDeleteRequest(w http.ResponseWriter, req *http.Request) {
	if rn.canAlterPeer() {
		var delReq cTypes.ConfigDeletionRequest

		if err := json.NewDecoder(req.Body).Decode(&delReq); err != nil {
			rn.writeError(w, http.StatusBadRequest, err)
		}

		confChange := &raftpb.ConfChange{
			NodeID: delReq.ID,
		}

		if err := rn.proposePeerDeletion(confChange, false); err != nil {
			rn.writeError(w, http.StatusInternalServerError, err)
		}

		rn.writeSuccess(w, nil)
	} else {
		rn.writeNodeNotReady(w)
	}
}

// wrapper to allow rn state to persist through handler func
func (rn *Node) peerAddHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		rn.handlePeerAddRequest(w, req)
	}
}

// if bootstrap node or in a cluster then accept these attempts,
// and wait for the message to be committed(err or retry after timeout?)
//
// Otherwise respond with an error that this node isn't in a state to add
// members
func (rn *Node) handlePeerAddRequest(w http.ResponseWriter, req *http.Request) {
	if rn.canAlterPeer() {
		var addReq cTypes.ConfigAdditionRequest

		if err := json.NewDecoder(req.Body).Decode(&addReq); err != nil {
			rn.writeError(w, http.StatusBadRequest, err)
		}

		var reqHost string
		if addReq.Host == "" {
			var err error
			reqHost, _, err = net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				rn.writeError(w, 500, err)
			}
		} else {
			reqHost = addReq.Host
		}
		confContext := cTypes.Peer{
			IP:                reqHost,
			RaftPort:          addReq.RaftPort,
			ConfigurationPort: addReq.ConfigurationPort,
		}

		confContextData, err := json.Marshal(confContext)
		if err != nil {
			rn.writeError(w, http.StatusInternalServerError, err)
		}

		confChange := &raftpb.ConfChange{
			NodeID:  addReq.ID,
			Context: confContextData,
		}

		if err := rn.proposePeerAddition(confChange, false); err != nil {
			rn.writeError(w, http.StatusInternalServerError, err)
		}

		addResp := &cTypes.ConfigAdditionResponseData{
			cTypes.ConfigPeerData{
				RaftPort:          rn.raftPort,
				ConfigurationPort: rn.configPort,
				ID:                rn.id,
				RemotePeers:       rn.peerMap,
			},
		}

		rn.writeSuccess(w, addResp)
	} else {
		rn.writeNodeNotReady(w)
	}
}

// TODO: Figure out how to handle these errs rather than just continue...
// thought of having a slice of accumulated errors?
// Or log.Warning on all failed attempts and if unsuccessful return a general failure
// error
func (rn *Node) requestRejoinCluster() error {
	var resp *http.Response
	var respData cTypes.ConfigServiceResponse

	if len(rn.bootstrapPeers) == 0 {
		return nil
	}

	for _, peer := range rn.bootstrapPeers {
		peerAPIURL := fmt.Sprintf("%s%s", peer, peerEndpoint)

		resp, err := http.Get(peerAPIURL)
		if err != nil {
			rn.logger.Warning(err.Error())
			continue
			//return err
		}

		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			rn.logger.Warning(err.Error())
			//return err
		}

		if respData.Status == peerServiceStatusError {
			continue
		} else if respData.Status == peerServiceStatusSuccess {

			var peerData cTypes.ConfigMembershipResponseData
			if err := json.Unmarshal(respData.Data, &peerData); err != nil {
				return errors.Wrap(err, "Error unmarshaling peer membership data")
			}

			return rn.addPeersFromRemote(peer, &peerData.ConfigPeerData)
		}
	}
	if respData.Status == peerServiceStatusError {
		return fmt.Errorf("Error %d - %s", resp.StatusCode, respData.Message)
	}
	// TODO: Should return the general error from here
	return errors.New("Couldn't connect to thingy")
}

func (rn *Node) addPeersFromRemote(remotePeer string, remoteMemberResponse *cTypes.ConfigPeerData) error {
	peerURL, err := url.Parse(remotePeer)
	if err != nil {
		return errors.Wrap(err, "Error parsing remote peer string for URL")
	}

	reqHost, _, err := net.SplitHostPort(peerURL.Host)
	if err != nil {
		return err
	}

	addURL := fmt.Sprintf("http://%s",
		net.JoinHostPort(reqHost, strconv.Itoa(remoteMemberResponse.RaftPort)))

	rn.transport.AddPeer(eTypes.ID(remoteMemberResponse.ID), []string{addURL})
	rn.logger.Info("Adding peer from HTTP request: %x\n", remoteMemberResponse.ID)
	rn.peerMap[remoteMemberResponse.ID] = cTypes.Peer{
		IP:                reqHost,
		RaftPort:          remoteMemberResponse.RaftPort,
		ConfigurationPort: remoteMemberResponse.ConfigurationPort,
	}
	rn.logger.Debugf("Current Peer Map: %v", rn.peerMap)

	for id, context := range remoteMemberResponse.RemotePeers {
		if id != rn.id {
			addURL := fmt.Sprintf("http://%s", net.JoinHostPort(context.IP, strconv.Itoa(context.RaftPort)))
			rn.transport.AddPeer(eTypes.ID(id), []string{addURL})
			rn.logger.Info("Adding peer from HTTP request: %x\n", id)
		}
		rn.peerMap[id] = context
		rn.logger.Debugf("Current Peer Map: %v", rn.peerMap)
	}
	return nil
}

func (rn *Node) requestSelfAddition() error {
	var resp *http.Response
	var respData cTypes.ConfigServiceResponse

	reqData := cTypes.ConfigAdditionRequest{
		ID:                rn.id,
		RaftPort:          rn.raftPort,
		ConfigurationPort: rn.configPort,
	}

	for _, peer := range rn.bootstrapPeers {
		mar, err := json.Marshal(reqData)
		if err != nil {
			rn.logger.Warning(err.Error())
			//return err
		}

		reader := bytes.NewReader(mar)
		peerAPIURL := fmt.Sprintf("%s%s", peer, peerEndpoint)

		resp, err = http.Post(peerAPIURL, "application/json", reader)
		if err != nil {
			rn.logger.Warning(err.Error())
			return err
		}

		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			rn.logger.Warning(err.Error())
			// return err
		}

		if respData.Status == peerServiceStatusError {
			continue
		} else if respData.Status == peerServiceStatusSuccess {

			// this ought to work since it should be added to cluster now
			var peerData cTypes.ConfigAdditionResponseData
			if err := json.Unmarshal(respData.Data, &peerData); err != nil {
				return errors.Wrap(err, "Error unmarshaling peer addition response")
			}

			return errors.Wrap(rn.addPeersFromRemote(peer, &peerData.ConfigPeerData), "Error add peer from remote data")
		}
	}
	if respData.Status == peerServiceStatusError {
		return fmt.Errorf("Error %d - %s", resp.StatusCode, respData.Message)
	}
	return errors.New("No available nodey thingy")
}

func (rn *Node) requestSelfDeletion() error {
	var resp *http.Response
	var respData cTypes.ConfigServiceResponse
	reqData := cTypes.ConfigDeletionRequest{
		ID: rn.id,
	}
	for id, peerData := range rn.peerMap {
		if id == rn.id {
			continue
		}
		mar, err := json.Marshal(reqData)
		if err != nil {
			return errors.Wrap(err, "Error marshalling peer deletion request")
		}

		reader := bytes.NewReader(mar)
		peerAPIURL := fmt.Sprintf("http://%s%s",
			net.JoinHostPort(peerData.IP, strconv.Itoa(peerData.ConfigurationPort)),
			peerEndpoint)

		req, err := http.NewRequest("DELETE", peerAPIURL, reader)
		if err != nil {
			return errors.Wrap(err, "Error creating new request for deleting peer")
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err = (&http.Client{}).Do(req)
		if err != nil {
			return errors.Wrap(err, "Error sending request to delete myself")
		}

		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return errors.Wrap(err, "Error decoding response for self deletion")
		}

		if respData.Status == peerServiceStatusSuccess {
			return nil
		}

	}
	if respData.Status == peerServiceStatusError {
		return fmt.Errorf("Error %d - %s", resp.StatusCode, respData.Message)
	}
	return nil
}

var peerServiceStatusSuccess = "success"
var peerServiceStatusError = "error"

var peerServiceNodeNotReady = "Invalid Node"

func (rn *Node) writeSuccess(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var respData []byte
	var err error
	if body != nil {
		respData, err = json.Marshal(body)
		if err != nil {
			rn.logger.Errorf(err.Error())
		}
	}

	if err = json.NewEncoder(w).Encode(cTypes.ConfigServiceResponse{Status: peerServiceStatusSuccess, Data: respData}); err != nil {
		rn.logger.Errorf(err.Error())
	}
}
func (rn *Node) writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(cTypes.ConfigServiceResponse{Status: peerServiceStatusError, Message: err.Error()}); err != nil {
		rn.logger.Errorf(err.Error())
	}
}

func (rn *Node) writeNodeNotReady(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(cTypes.ConfigServiceResponse{Status: peerServiceStatusError, Message: peerServiceNodeNotReady}); err != nil {
		rn.logger.Errorf(err.Error())
	}
}
