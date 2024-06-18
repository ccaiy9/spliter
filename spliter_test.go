package spliter

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func Test_correctBatchSqlCorrect(t *testing.T) {

	sqlTestList := []string{
		"INSERT INTO `yu`.table_name(`insert into`, `cd`) values (1,7),(`中文`,`shanghai`),(`普通话`,`（中国上海)`),(6,7);",
		"insert into  `test_sb`(ty, trans, pos) VALUES (`USA CHina 北京`,`  `, `VALUES`), (`熊猫`, position, `(الفنون الجميلة)`), (`;`, `fine arts`, pety);",
		"insert INTO  test_sb(`values`, `VALUES`, `pos`) VALUES    (`insert into`, `select * from other`, `VALUES`), (`(yu)`, `)`,`;`);",
		"INSERT   INTO   yu.`table_name`   values (1,3),(w,p6),(wp,1);",
		"INSERT   INTo  `yu`.`table_name` vaLues (1257,34),(`ww`,`po`),(`wp`,`1`);",
		`INSErT   InTO  yu.table_name values (wui,yu),("   689","899\n"),   ("wtux211 \t ","09sxnaj"), (1,2);`,
		`iNSERT   INTO  yu.table_name Values (wui,yu),("wempohu\\"),("689","899\n"),
		                ("wtux211 \t ","09sxnaj");`,
		`INSERT
		INTO  yu.table_name values (wu,opx),("天才","bbo天\n"),
		                           ("烫烫烫\t ","09sxnaj");`,

		`INSERT
		INTO
		yu.table_name
		values (wu,opx),
		("烫烫烫","上海省\n"),
		("共opp \t ","09sxnaj");`,

		"REPLACE INTO `yu`.table_name(`REPLACE into`, `cd`) values (12,34),(`ww`,po),(wp,`1)`);",
		"replace into yu.table_name(wuyou, t2) values (`...`, `  %^&#@!!`),(7,2);",
	}

	Option := &Option{
		splitCompeteSqlChunk: 10,
	}

	for _, sql := range sqlTestList {
		fmt.Println("【begin】-------- sql --------")
		stms, err := SplitBatchInsertSql(&sql, Option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, st := range stms {
			fmt.Println(st)
		}

		fmt.Println("【end】-------- sql ----------------")
	}
}

func Test_incorrectBatchSqlCorrect(t *testing.T) {
	sqlTestList := []string{
		"INSER   INTO  yu.table_name values (12,34),(`ww`,`po`),(`wp`,`1`);",
		"INSERt INTO  table_name values 12,34),(`ww`,`po`),(`wp`,`1`);",
		"INSERt INTO  table_name values (12,34),(`ww`,`po`);(`wp`,`1`);",
		"INSERt INTO  table_name values (12,34),(`ww`,`po`),(`wp`,,`1`);",
	}

	for index, sql := range sqlTestList {
		fmt.Println(index)
		Option := &Option{
			AbsPath:              "/Users/Downloads/spliter/",
			Index:                strconv.Itoa(index),
			splitCompeteSqlChunk: 10,
		}

		fmt.Println("【begin】-------- sql --------")
		stms, err := SplitBatchInsertSql(&sql, Option)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, st := range stms {
			fmt.Println(st)
		}

		fmt.Println("【end】-------- sql ----------------")

	}

}
