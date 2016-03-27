package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

const PORT string = ":8081"
const CA_MSG string = "supersecretCAhashmessage"

// System Paramters
const sys_p string = "168199388701209853920129085113302407023173962717160229197318545484823101018386724351964316301278642143567435810448472465887143222934545154943005714265124445244247988777471773193847131514083030740407543233616696550197643519458134465700691569680905568000063025830089599260400096259430726498683087138415465107499"
const sys_q string = "959452661475451209325433595634941112150003865821"
const sys_g string = "94389192776327398589845326980349814526433869093412782345430946059206568804005181600855825906142967271872548375877738949875812540433223444968461350789461385043775029963900638123183435133537262152973355498432995364505138912569755859623649866375135353179362670798771770711847430626954864269888988371113567502852"

// CA private and public key pair
const SK_CA string = "432398415306986194693973996870836079581453988813"
const PK_CA string = "49336018324808093534733548840411752485726058527829630668967480568854756416567496216294919051910148686186622706869702321664465094703247368646506821015290302480990450130280616929226917246255147063292301724297680683401258636182185599124131170077548450754294083728885075516985144944984920010138492897272069257160"

func generateCert(connection net.Conn, user string, userPubKey *big.Int,
	dss_r *big.Int, dss_s *big.Int, dss_hash *big.Int, expDate time.Time) {

	fmt.Printf("dss_r = %v \ndss_s = %v\n", dss_r, dss_s)
	fmt.Printf("expiry date = %v\n", expDate)

	connection.Write([]byte("test" + "\n"))
}

func main() {
	fmt.Println("Launching the goca server...")

	ln, _ := net.Listen("tcp", PORT)

	connection, _ := ln.Accept()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Listening to goca-client request on port: ", string(PORT))

	user, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Println("User identity requesting certificate: ", string(user))

	userPubKeyStr, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Println("Public Key of the connected user: ", string(userPubKeyStr))
	userPubKey := *big.NewInt(0)
	userPubKey.SetString(userPubKeyStr, 10)

	fmt.Print("Enter the expiry date for this certificate (yyyy-mmm-dd): ")
	date, _ := reader.ReadString('\n')
	const shortForm = "2006-Jan-02"
	expDate, _ := time.Parse(shortForm, date)

	// System Parameters
	dss_p := *big.NewInt(0)
	dss_p.SetString(sys_p, 10)

	dss_q := *big.NewInt(0)
	dss_q.SetString(sys_q, 10)

	dss_g := *big.NewInt(0)
	dss_g.SetString(sys_g, 10)

	dss_x := *big.NewInt(0)
	dss_x.SetString(SK_CA, 10)

	dss_hash := md5.New()

	var dss_k *big.Int
	dss_k, _ = rand.Int(rand.Reader, &dss_q)

	dss_r := *big.NewInt(0)
	dss_r.Mod((dss_g.Exp(&dss_g, dss_k, &dss_p)), &dss_q)

	dss_s := big.NewInt(0)
	io.WriteString(dss_hash, CA_MSG)
	dss_hash_bigInt := big.NewInt(0)
	dss_hash_bigInt.SetString(string(dss_hash.Sum(nil)), 10)

	dss_s = dss_k.ModInverse(dss_k.Add(dss_hash_bigInt, dss_x.Mul(&dss_x, &dss_r)), &dss_q)

	generateCert(connection, user, &userPubKey, dss_s, &dss_r, dss_hash_bigInt, expDate)
}
