Model holds data, no ui, no logic

View reads Model and produces UI

func(m model) View() string {

	return m.msg
}

so view CAN READ       the model
        CANNOT CHANGE  the model
        CAN SEND messages
                    messages represent events not actions

UPDATE is the only place where model can change 

take a msg
take current model
return new model

func (m model) Update(msg tea.Msg) (tea.Model,tea.Cmd){

	return m,nil
}

Init creates the initial model

Init -> Model
Model -> View
View -> Msg
Msg -> Update
Update -> Model
Model -> View
and so on