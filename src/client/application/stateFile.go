package application

import (
	"context"
	"fmt"
	"os"

	"gophKeeper/src/client/cfg"
	GRPCClient "gophKeeper/src/client/grpcClient"
)

type stateFileType struct {
	d *stateDataType
}

func (s *stateFileType) getName() string {
	return s.d.getName()
}

func newFileState(app *Application, config *cfg.ConfigT) state {
	st := &stateDataType{
		stateGetName: stateGetName{stateName: "File view"},
		data:         make(dataType),
		dataIDs:      make([]string, 0),
		app:          app,
		category:     GRPCClient.CategoryFile,
		fields: []string{
			"filename",
			"metadata",
		},
		config:     config,
		inputField: -1,
		fn: &accessFunctions{
			get: app.grpc.LoadFileData,
			add: app.grpc.StoreFileData,
		},
	}
	fileType := &stateFileType{d: st}
	st.local = &localFunctions{
		load: fileType.load,
		add:  fileType.add,
	}
	return fileType
}

func (s *stateFileType) execute(ctx context.Context, command string) (resultState state, err error) {
	res, err := s.d.execute(ctx, command)
	switch res := res.(type) {
	case *stateDataType:
		s.d = res
		return s, err
	default:
		return res, err
	}
}

func (s *stateFileType) add(ctx context.Context, command string) (resultState state, err error) {
	if s.d.inputField != len(s.d.fields)-1 {
		s.d.currentInput = append(s.d.currentInput, command)
		s.d.inputField++
		fmt.Printf("Input %q: ", s.d.fields[s.d.inputField])
		return s, nil
	}

	file, err := os.ReadFile(s.d.currentInput[0])
	if err != nil {
		return s, fmt.Errorf("file open failed: %w", err)
	}
	data := s.d.app.encoder.Encode(string(file))
	dataID, metadata, err := s.d.fn.add(ctx, data, command)
	s.d.inputField = -1
	if err != nil {
		return s, fmt.Errorf("adding entry failed: %w", err)
	}
	s.d.data[dataID] = &dataEntry{DataID: dataID, Metadata: metadata, Data: s.d.currentInput}
	s.d.currentInput = []string{}
	return s, nil
}

func (s *stateFileType) load(ctx context.Context, id int) (resultState state, err error) {
	dataID := s.d.dataIDs[id]
	if len(s.d.data[dataID].Data) == 0 {
		var resp []byte
		resp, err = s.d.fn.get(ctx, s.d.data[dataID].DataID)
		if err != nil {
			return s, fmt.Errorf("requesting data failed: %w", err)
		}

		var data string
		data, err = s.d.app.encoder.Decode(resp)
		if err != nil {
			return s, fmt.Errorf("decoding data failed: %w", err)
		}
		err = os.WriteFile(dataID, []byte(data), 0o600)
		if err != nil {
			return s, fmt.Errorf("writng file failed: %w", err)
		}
		s.d.data[dataID].Data = []string{dataID}
	}
	fmt.Println("file saved as ", dataID)
	fmt.Println("metadata: ", s.d.data[dataID].Metadata)
	return s, nil
}
