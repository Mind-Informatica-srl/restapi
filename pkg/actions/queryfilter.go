package actions

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var sqlOperators = map[string]string{
	"equal":          "=",
	"notequal":       "<>",
	"equalnumber":    "=",
	"notequalnumber": "<>",
	"equaldate":      "=",
	"notequaldate":   "<>",
	"equalboolean":   "=",
	"like":           "like",
	"likestart":      "like",
	"likeend":        "like",
	"gt":             ">",
	"lt":             "<",
	"gte":            ">=",
	"lte":            "<=",
	"gtdate":         ">",
	"ltdate":         "<",
	"gtedate":        ">=",
	"ltedate":        "<=",
	"isnull":         "is null",
	"isnotnull":      "is not null",
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var onlyNumber = regexp.MustCompile(`[0-9\.\,]*`)

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func QueryFilter(db *gorm.DB, url *url.URL) (*gorm.DB, error) {
	var q string
	if params, ok := url.Query()["q"]; ok {
		q = params[0]
		if q != "" {
			//se c'è un criterio di ricerca
			//si spitta per &
			searchCriteria := strings.Split(q, "&")
			for _, c := range searchCriteria {
				//criterio è della forma [nomeAttributo].[operatore]=[valore] oppure [nome_tabella1]|[nomeAttributoTabella1].[...].[nome_tabellaN]|[nomeAttributoTabellaN].[nomeAttributo].[operatore]=[valore]
				//per ogni criterio si splitta per il =
				criteria := strings.Split(c, "=") //la parte sinistra rappresenta il campo e l'operatore, la parte di destra il valore
				if len(criteria) < 2 {
					return nil, errors.New("parametri di ricerca non validi")
				}
				//si splitta la parte sinistra per il .
				p := strings.Split(criteria[0], ".")
				valore := criteria[1]
				condition, values, err := prepareSubQuery(p, valore)
				print(condition)
				if err != nil {
					return nil, err
				}
				if values != nil {
					db = db.Where(condition, values...)
				} else {
					// values è nil qualora l'operatore è isnull o isnotnull
					db = db.Where(condition)
				}
				/*if len(p) > 2 {
					//caso complesso
					condition := prepareSubQuery(db, p, valore)
					db = db.Where(condition)
					//clienti|ClienteID.Localita.like=FIR
					//db.Where("cliente_id in (?)", db.Table("clienti").Select("id").Where("lower(localita) like lower('%FIR%'")))
					// SELECT * FROM "pratiche" WHERE cliente_id in (SELECT id FROM "clienti" where lower(localita) like lower('%FIR%'));

					//subQuery := db.Select("AVG(age)").Where("name LIKE ?", "name%").Table("users")
					//db.Select("AVG(age) as avgage").Group("name").Having("AVG(age) > (?)", subQuery).Find(&results)
					// SELECT AVG(age) as avgage FROM `users` GROUP BY `name` HAVING AVG(age) > (SELECT AVG(age) FROM `users` WHERE name LIKE "name%")
				} else {
					//caso [nomeAttributo].[operatore]=[valore]
					condition := getWhereCondition(p[0], p[1], valore)
					db = db.Where(condition)
				}*/
			}
		}
	}
	return db, nil
}

func composeCriteria(attributeName string, operatorName string, stringValue string) (string, string, string) {
	//c := Cliente{}
	//field, ok := reflect.TypeOf(c).Elem().FieldByName("nome attributo")
	//tag = string(field.Tag)
	columnName := ToSnakeCase(attributeName)
	operator := sqlOperators[operatorName]
	onlyNumber.FindString(stringValue)
	var value string
	switch operatorName {
	case "equal":
		if onlyNumber.FindString(stringValue) == stringValue {
			value = stringValue
		} else {
			value = strings.ToLower(stringValue)
			columnName = "lower(" + columnName + ")"
		}
	case "notequal":
		value = strings.ToLower(stringValue)
		columnName = "lower(" + columnName + ")"
	case "equalboolean":
		value = "'" + stringValue + "'"
	case "like":
		value = "%" + strings.ToLower(stringValue) + "%"
		columnName = "lower(" + columnName + ")"
	case "likestart":
		value = "%" + strings.ToLower(stringValue) + "%"
		columnName = "lower(" + columnName + ")"
	case "likeend":
		value = "%" + strings.ToLower(stringValue) + "%"
		columnName = "lower(" + columnName + ")"
	case "isnull":
		value = ""
	case "isnotnull":
		value = ""
	default:
		value = stringValue
	}
	if strings.HasSuffix(operatorName, "date") {
		columnName = columnName + "::date"
		value = "'" + strings.Split(value, "T")[0] + "'"
	}
	return columnName, operator, value
}

func getWhereCondition(attributeName string, operatorName string, stringValue string) (string, []interface{}, error) {
	if strings.Contains(attributeName, "#") {
		cols := strings.Split(attributeName, "#")
		values := strings.Split(stringValue, "#")
		vs := make([]interface{}, len(values))
		if len(cols) != len(values) {
			return "", []interface{}{}, errors.New("filtro di ricerca non corretto")
		}
		var condition strings.Builder
		var writeAnd = false
		for i := 0; i < len(cols); i++ {
			c, o, v := composeCriteria(cols[i], operatorName, values[i])
			vs[i] = v
			if writeAnd {
				condition.WriteString(" and ")
			}
			condition.WriteString("(")
			condition.WriteString(c)
			condition.WriteString(" ")
			condition.WriteString(o)
			condition.WriteString(" ")
			condition.WriteString("?")
			condition.WriteString(")")
			writeAnd = true
		}
		return condition.String(), vs, nil
	} else {
		c, o, v := composeCriteria(attributeName, operatorName, stringValue)
		if strings.HasSuffix(operatorName, "null") {
			return c + " " + o, nil, nil
		}
		return c + " " + o + " ?", []interface{}{v}, nil
	}
}

// Esempio: select * from notule where FatturaProgressivo in (
// 	select Progressivo from fatture where clienteID in (
// 		select ClienteID from referenti where Nominativo = 'Marco'
// ))
//FatturaProgressivo|fatture|Progressivo.ClienteID|referenti_clienti|ClienteID.Nominativo.equal
func prepareSubQuery(p []string, stringValue string) (string, []interface{}, error) {
	l := len(p)
	var operatore, parametro string
	var whereCondition string
	var values []interface{}
	var err error
	for i := l - 1; i >= 0; i-- {
		switch i {
		case l - 1:
			operatore = p[l-1] //equal
		case l - 2:
			parametro = p[l-2]                                                                 //Nominativo
			whereCondition, values, err = getWhereCondition(parametro, operatore, stringValue) //Nominativo = 'Marco'
			if err != nil {
				return "", []interface{}{}, errors.New("filtro non corretto")
			}
		default:
			a := strings.Split(p[i], "|") //ClienteID|referenti_clienti|ClienteID
			if len(a) == 3 {
				whereCondition = ToSnakeCase(a[0]) + " in (select " + ToSnakeCase(a[2]) + " from " + a[1] + " where (" + whereCondition + ") )"
			} else {
				return "", []interface{}{}, errors.New("filtro non corretto")
			}
		}
	}
	return whereCondition, values, nil
}

//si cerca in selectColumns se si deve fare il preload di altre tabelle
func QueryPreload(db *gorm.DB, selectColumns []string) *gorm.DB {
	if selectColumns != nil {
		//columnNames := ""
		for i := 0; i < len(selectColumns); i++ {
			if strings.Contains(selectColumns[i], ".") {
				//se l'i-esimo selectColumns ha un ".", si splitta e si prende il primo valore per il preload (es: Cliente.RagioneSociale => si fa il preload di Cliente)
				ref := strings.Split(selectColumns[i], ".")[0]
				db = db.Preload(ref)
			}
		}
		//db = db.Select(columnNames)
	}
	return db
}

//NON VA BENE SE CI SONO CAMPI CALCOLATI
//si cerca in selectColumns l'elenco delle colonne su cui fare la select
/*func QueryPreload(db *gorm.DB, selectColumns []string) *gorm.DB {
	if selectColumns != nil {
		columnNames := ""
		for i := 0; i < len(selectColumns); i++ {
			if strings.Contains(selectColumns[i], ".") {
				//se l'i-esimo selectColumns ha un ".", si splitta e si prende il primo valore per il preload (es: Cliente.RagioneSociale => si fa il preload di Cliente)
				ref := strings.Split(selectColumns[i], ".")[0]
				db = db.Preload(ref)
			} else {
				//altrimenti si aggiunge il nome della colonna in snakecase a columnNames (es: PartitaIva => partita_iva)
				if columnNames != "" {
					columnNames = columnNames + ","
				}
				columnNames = columnNames + ToSnakeCase(selectColumns[i])
			}
		}
		db = db.Select(columnNames)
	}
	return db
}
*/
