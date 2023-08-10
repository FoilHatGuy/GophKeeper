package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gophKeeper/src/client/cfg"
	GRPCClient "gophKeeper/src/client/grpcClient"
)

type dataEntry struct {
	DataID   string
	Metadata string
	Data     []string
}
type (
	dataType       map[string]*dataEntry
	stateCredsType struct {
		stateName    string
		app          *Application
		config       *cfg.ConfigT
		data         dataType
		dataIDs      []string
		inputField   int
		currentInput []string
	}
)

var (
	fieldsCred  = []string{"login   ", "password", "metadata"}
	commandAdd  = []string{"add", "a"}
	commandHead = []string{"head", "h"}
)

func newCredsState(app *Application, config *cfg.ConfigT) state {
	return &stateCredsType{
		app:        app,
		config:     config,
		inputField: -1,
		stateName:  "Credentials view",
	}
}

func (s *stateCredsType) execute(ctx context.Context, command string) (resultState state, err error) {
	arguments := strings.Split(command, " ")
	if s.inputField >= 0 {
		return s.add(ctx, command)
	}

	switch {
	case includes(commandHelp, strings.ToLower(arguments[0])):
		fmt.Printf("This is Login screen. Available commands:\n"+
			"help - %q - shows available commands (this screen)\n"+
			"load - %q $id$ - requests encoded data for entity:\n"+
			"add  - %q - initiates adding of data:\n"+
			"head - %q - shows saved entries:\n"+
			"back - %q - return to menu",
			commandHelp,
			commandLoad,
			commandAdd,
			commandHead,
			commandBack)

	case includes(commandAdd, strings.ToLower(arguments[0])):
		s.inputField = 0
		fmt.Printf("Input %q: ", fieldsCred[s.inputField])
		return s, nil

	case includes(commandHead, strings.ToLower(arguments[0])):
		err = s.fetch(ctx)
		if err != nil {
			return nil, fmt.Errorf("head fetching failed: %w", err)
		}
		if len(s.data) == 0 {
			return s, errors.New("empty list")
		}
		s.show(ctx)
		return s, nil

	case includes(commandLoad, strings.ToLower(arguments[0])):
		if len(arguments) != 2 {
			return s, fmt.Errorf("%w\nyou should specify the id of entry", ErrUnrecognizedCommand)
		}
		var id int
		id, err = strconv.Atoi(arguments[1])
		if err != nil {
			return s, fmt.Errorf("%w\nid should be int", ErrUnrecognizedCommand)
		}
		dataID := s.dataIDs[id]
		if len(s.data[dataID].Data) == 0 {
			var resp []byte
			resp, err = s.app.grpc.LoadCredentials(ctx, s.data[dataID].DataID)
			if err != nil {
				return s, fmt.Errorf("requesting data failed: %w", err)
			}

			var data string
			data, err = s.app.encoder.Decode(resp)
			if err != nil {
				return s, fmt.Errorf("decoding data failed: %w", err)
			}
			dataArr := strings.Split(data, "\x00")
			s.data[dataID].Data = dataArr
		}
		for i := 0; i < len(fieldsCred)-1; i++ {
			fmt.Println(fieldsCred[i]+": ", s.data[dataID].Data[i])
		}
		fmt.Println("metadata: ", s.data[dataID].Metadata)

		return s, nil

	case includes(commandBack, strings.ToLower(arguments[0])):
		return s.app.cat[stateMenu], nil

	default:
		return s, ErrUnrecognizedCommand
	}
	s.show(ctx)
	return s, nil
}

func (s *stateCredsType) add(ctx context.Context, command string) (resultState state, err error) {
	if s.inputField != len(fieldsCred)-1 {
		s.currentInput = append(s.currentInput, command)
		s.inputField++
		fmt.Printf("Input %q: ", fieldsCred[s.inputField])
		return s, nil
	}

	data := s.app.encoder.Encode(strings.Join(s.currentInput, "\x00"))
	dataID, metadata, err := s.app.grpc.StoreCredentials(ctx, data, command)
	s.inputField = -1
	if err != nil {
		return s, fmt.Errorf("adding entry failed: %w", err)
	}
	s.data[dataID] = &dataEntry{DataID: dataID, Metadata: metadata}
	return s, nil
}

func (s *stateCredsType) show(_ context.Context) {
	const metaLen = 50
	const dataLen = 30
	for i, dataID := range s.dataIDs {
		el := s.data[dataID]
		fmt.Println(
			fmt.Sprintf("%*d", 5, i),
			fmt.Sprintf("%-*q", metaLen+2, firstN(el.Metadata, metaLen)),
			fmt.Sprintf("%-*q", dataLen+2, firstN(strings.Join(el.Data, " "), dataLen)),
		)
	}
}

func (s *stateCredsType) fetch(ctx context.Context) (err error) {
	if s.data == nil {
		s.data = make(dataType)
	}
	head, err := s.app.grpc.GetCategoryHead(ctx, GRPCClient.CategoryCred)
	if err != nil {
		return fmt.Errorf("during fetching category head: %w", err)
	}

	for _, el := range head {
		if _, ok := s.data[el.DataID]; !ok {
			s.data[el.DataID] = &dataEntry{
				DataID:   el.DataID,
				Metadata: el.Metadata,
			}
			s.dataIDs = append(s.dataIDs, el.DataID)
		}
	}
	return nil
}

func (s *stateCredsType) getName() string {
	return s.stateName
}
