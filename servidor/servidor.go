package servidor

import (
	"crud/banco"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type usuario struct {
	ID    uint32 `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

//criar usuario
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpo, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario
	if err = json.Unmarshal(corpo, &usuario); err != nil {
		w.Write([]byte("erro ao converter usuario"))
		return
	}
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha conectarf no banco"))
		return
	}
	defer db.Close()
	statement, err := db.Prepare("insert into usuarios (nome, email) values (?,?)")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
		return
	}
	defer statement.Close()
	insercao, err := statement.Exec(usuario.Nome, usuario.Email)
	if err != nil {
		w.Write([]byte("erro ao executar o statement"))
		return
	}

	idInserido, err := insercao.LastInsertId()
	if err != nil {
		w.Write([]byte("erro ao obter o id inserido"))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("Usuario insirido com sucesso! id: %d", idInserido)))
}

// busca todos os usuarios
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("erro ao conectar com o banco de dados"))
		return
	}

	defer db.Close()
	linhas, err := db.Query("select * from usuarios")
	if err != nil {
		w.Write([]byte("erro ao listar usuarios"))
		return
	}
	defer linhas.Close()
	var usuarios []usuario
	for linhas.Next() {
		var usuario usuario
		if err := linhas.Scan(&usuario.ID, &usuario.Nome, &usuario.Email); err != nil {
			w.Write([]byte("erro ao scanear usuarios"))
			return
		}
		usuarios = append(usuarios, usuario)
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("erro ao converter usuarios para json"))
		return
	}
}

// busca um usuario expecifico
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	parametros := mux.Vars(r)
	ID, err := strconv.ParseUint(parametros["id"], 10, 32)
	if err != nil {
		w.Write([]byte("erro ao converter ID"))
		return
	}
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("erro ao conectar com o banco de dados"))
		return
	}

	defer db.Close()
	linha, err := db.Query("select * from usuarios where id = ?", ID)
	if err != nil {
		w.Write([]byte("erro ao listar usuarios"))
		return
	}
	defer linha.Close()
	var usuarios usuario
	if linha.Next() {
		if err := linha.Scan(&usuarios.ID, &usuarios.Nome, &usuarios.Email); err != nil {
			w.Write([]byte("erro ao converter usuarop"))

		}
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(usuarios); err != nil {
		w.Write([]byte("erro ao converter usuarios para json"))
		return
	}
}

//atualiza um usuario
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	corpo, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var usuario usuario
	if err = json.Unmarshal(corpo, &usuario); err != nil {
		w.Write([]byte("erro ao converter usuario"))
		return
	}
	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha conectar no banco"))
		return
	}
	defer db.Close()

	linha, err := db.Query("select * from usuarios where id = ?", usuario.ID)
	if err != nil {
		w.Write([]byte("erro ao listar usuarios"))
		return
	}
	if !linha.Next() {
		w.Write([]byte("id invalido de usuario"))
		return
	}
	statement, err := db.Prepare("update usuarios set nome =? , email=? where id=?")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
		return
	}
	defer statement.Close()
	if _, err := statement.Exec(usuario.Nome, usuario.Email, usuario.ID); err != nil {
		w.Write([]byte("erro ao executar o statement"))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("Usuario atualizado com sucesso! id: %d", usuario.ID)))
}

//remove um usuario
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	ID, err := strconv.ParseUint(parametros["id"], 10, 32)
	if err != nil {
		w.Write([]byte("erro ao converter ID"))
		return
	}

	db, err := banco.Conectar()
	if err != nil {
		w.Write([]byte("Falha conectar no banco"))
		return
	}
	defer db.Close()

	linha, err := db.Query("select * from usuarios where id = ?", ID)
	if err != nil {
		w.Write([]byte("erro ao listar usuarios"))
		return
	}
	if !linha.Next() {
		w.Write([]byte("id invalido de usuario"))
		return
	}
	statement, err := db.Prepare("delete from usuarios  where id=?")
	if err != nil {
		w.Write([]byte("erro ao criar o statement"))
		return
	}
	defer statement.Close()
	retorno, err := statement.Exec(ID)
	if err != nil {
		w.Write([]byte("erro ao executar o statement"))
		return
	}
	if _, err := retorno.RowsAffected(); err != nil {
		w.Write([]byte("erro ao obter as linhas afetadas"))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf("Usuario removido com sucesso! id: %d", ID)))
}
