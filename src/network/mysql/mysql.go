package mysql
import(
    "database/sql"

    _ "github.com/go-sql-driver/mysql"
)

func GetConnection()(ret *sql.DB, err error){
    ret, err = sql.Open("mysql", "root:@/code_golf")
    return
}
