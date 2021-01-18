package war

type WarProviderInterface interface {
	GetWars() ([]War, error)
	PutWars(wars []War) error
}

type WarProviderSimpleList struct {
	wars []War
}

func NewWarProviderSimpleList() WarProviderSimpleList {
	return WarProviderSimpleList{}
}

func (simpleList WarProviderSimpleList) GetWars() ([]War, error) {
	return simpleList.wars, nil
}

func (simpleList *WarProviderSimpleList) PutWars(warsToAdd []War) error {

	for _, warToAdd := range warsToAdd {
		foundExistingWar := false
		for existingWarIndex, existingWar := range simpleList.wars {
			if existingWar.Name == warToAdd.Name {
				simpleList.wars[existingWarIndex] = warToAdd
				foundExistingWar = true
				break
			}
		}
		if !foundExistingWar {
			simpleList.wars = append(simpleList.wars, warToAdd)
		}
	}

	return nil
}

var warProviderSimpleListInterfaceChecker WarProviderInterface = &WarProviderSimpleList{}
