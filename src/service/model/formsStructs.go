package structs

//Field usada para guardar os dados para servirem de base para construção dinâmica dos campos em tela
type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Title string `json:"title"`
}

//Form usada para guardar os dados para servirem de base para construção dinâmica de forms em tela
type Form struct {
	FormName string  `json:"formName"`
	Fields   []Field `json:"fields"`
}