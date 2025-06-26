package services

import (
	"database/sql"
	//"fmt"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/segmentio/kafka-go"
	"time"
)

type DBController struct {
	DB *sql.DB
}

func NewDBController(filepath string) (*DBController, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	return &DBController{DB: db}, nil
}

// ---------- Tabela de sessões ----------
type SessaoCliente struct {
	ID        int
	UserID    string
	JwtToken  string
	CreatedAt time.Time
	LastSeen  time.Time
}

func (dbc *DBController) InitSessaoTable() error {
	stmt := `
	CREATE TABLE IF NOT EXISTS sessao_cliente (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		jwt_token TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		last_seen DATETIME NOT NULL
	)`
	_, err := dbc.DB.Exec(stmt)
	return err
}

func (dbc *DBController) UpsertSessao(userID, jwt string) error {
	now := time.Now()
	var exists bool
	err := dbc.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM sessao_cliente WHERE user_id = ?)", userID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = dbc.DB.Exec(`UPDATE sessao_cliente SET jwt_token = ?, last_seen = ? WHERE user_id = ?`, jwt, now, userID)
	} else {
		_, err = dbc.DB.Exec(`INSERT INTO sessao_cliente (user_id, jwt_token, created_at, last_seen) VALUES (?, ?, ?, ?)`, userID, jwt, now, now)
	}
	return err
}

func (dbc *DBController) GetSessaoByUserID(userID string) (*SessaoCliente, error) {
	row := dbc.DB.QueryRow(`SELECT id, user_id, jwt_token, created_at, last_seen FROM sessao_cliente WHERE user_id = ?`, userID)
	var s SessaoCliente
	err := row.Scan(&s.ID, &s.UserID, &s.JwtToken, &s.CreatedAt, &s.LastSeen)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// ---------- Tabela central ----------
type Central struct {
	ID              int
	Identificador   string
	Nome            string
	Versao          string
	DataAtualizacao time.Time
}

func (dbc *DBController) InsertCentral(c Central) (int64, error) {
	stmt := `INSERT INTO central (identificador, nome, versao, data_atualizacao) VALUES (?, ?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, c.Identificador, c.Nome, c.Versao, c.DataAtualizacao)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dbc *DBController) GetCentralByID(id int) (*Central, error) {
	row := dbc.DB.QueryRow(`SELECT id, identificador, nome, versao, data_atualizacao FROM central WHERE id = ?`, id)
	var c Central
	err := row.Scan(&c.ID, &c.Identificador, &c.Nome, &c.Versao, &c.DataAtualizacao)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (dbc *DBController) UpdateCentral(c Central) error {
	stmt := `UPDATE central SET identificador = ?, nome = ?, versao = ?, data_atualizacao = ? WHERE id = ?`
	_, err := dbc.DB.Exec(stmt, c.Identificador, c.Nome, c.Versao, c.DataAtualizacao, c.ID)
	return err
}

func (dbc *DBController) DeleteCentral(id int) error {
	_, err := dbc.DB.Exec(`DELETE FROM central WHERE id = ?`, id)
	return err
}

// ---------- Tabela micro_controlador ----------
type MicroControlador struct {
	ID             string
	Tipo           string
	CaminhoBinario string
	Token          string
	IP             string
}

func (dbc *DBController) InsertMicro(m MicroControlador) (int64, error) {
	stmt := `INSERT INTO micro_controlador (tipo_controlador, caminho_binario, token, ip) VALUES (?, ?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, m.Tipo, m.CaminhoBinario, m.Token, m.IP)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dbc *DBController) GetMicroByID(id int) (*MicroControlador, error) {
	row := dbc.DB.QueryRow(`SELECT id, tipo_controlador, caminho_binario, token, ip FROM micro_controlador WHERE id = ?`, id)
	var m MicroControlador
	err := row.Scan(&m.ID, &m.Tipo, &m.CaminhoBinario, &m.Token, &m.IP)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (dbc *DBController) UpdateMicro(m MicroControlador) error {
	stmt := `UPDATE micro_controlador SET tipo_controlador = ?, caminho_binario = ?, token = ?, ip = ? WHERE id = ?`
	_, err := dbc.DB.Exec(stmt, m.Tipo, m.CaminhoBinario, m.Token, m.IP, m.ID)
	return err
}

func (dbc *DBController) DeleteMicro(id int) error {
	_, err := dbc.DB.Exec(`DELETE FROM micro_controlador WHERE id = ?`, id)
	return err
}

// ---------- Tabela topicos ----------
type Topico struct {
	ID        int
	Nome      string
	Mensagem  string
	CentralID int
}

func (dbc *DBController) InsertTopico(t Topico) (int64, error) {
	stmt := `INSERT INTO topicos (nome, mensagem, central_id) VALUES (?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, t.Nome, t.Mensagem, t.CentralID)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dbc *DBController) GetTopicoByID(id int) (*Topico, error) {
	row := dbc.DB.QueryRow(`SELECT id, nome, mensagem, central_id FROM topicos WHERE id = ?`, id)
	var t Topico
	err := row.Scan(&t.ID, &t.Nome, &t.Mensagem, &t.CentralID)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (dbc *DBController) UpdateTopico(t Topico) error {
	stmt := `UPDATE topicos SET nome = ?, mensagem = ?, central_id = ? WHERE id = ?`
	_, err := dbc.DB.Exec(stmt, t.Nome, t.Mensagem, t.CentralID, t.ID)
	return err
}

func (dbc *DBController) DeleteTopico(id int) error {
	_, err := dbc.DB.Exec(`DELETE FROM topicos WHERE id = ?`, id)
	return err
}

// ---------- Tabela treino ----------
type Treino struct {
	ID                int
	MicroID           string
	UserID            string
	PlanID            string
	MainMuscles       string
	ExerciseID        string
	GymID             string
	Summary           string
	WorkoutScore      float64
	ExerciseScore     float64
	RestingTime       float64
	UsedLoad          float64
	FailedReps        int
	TotalReps         int
	TotalSeries       int
	AmplitudeScore    float64
	DetailedAmplitude float64
	TimeScore         float64
	DetailedTimeScore float64
	JWT               string
}

func (dbc *DBController) InsertTreino(t Treino) (int64, error) {
	now := time.Now()
	stmt := `INSERT INTO treino (micro_id, user_id, failed_reps, total_reps, total_series, data_atualizacao) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, t.MicroID, t.UserID, t.FailedReps, t.TotalReps, t.TotalSeries, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Get, Update, Delete para treino podem ser adicionados da mesma forma conforme necessário

// ---------- Tabela versao_online ----------
type VersaoOnline struct {
	ID              int
	CentralID       string
	CaminhoBinario  string
	Versao          string
	DataVersao      time.Time
	TipoControlador string
}

func (dbc *DBController) InsertVersaoOnline(v VersaoOnline) (int64, error) {
	stmt := `INSERT INTO versao_online (central_id, caminho_binario, versao, data_versao, tipo_controlador) VALUES (?, ?, ?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, v.CentralID, v.CaminhoBinario, v.Versao, v.DataVersao, v.TipoControlador)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dbc *DBController) GetVersaoOnlineByMicroType(tipo_controlador string) (*VersaoOnline, error) {
	stmt := `SELECT versao, central_id, caminho_binario, tipo_controlador FROM versao_online WHERE tipo_controlador = ?`
	row := dbc.DB.QueryRow(stmt, tipo_controlador)
	var cm VersaoOnline
	err := row.Scan(&cm.Versao, &cm.CentralID, &cm.CaminhoBinario, &cm.TipoControlador)
	if err != nil {
		return nil, err
	}
	return &cm, nil
}

// ---------- Tabela central_micro ----------
type CentralMicro struct {
	Versao          string
	DataAtualizacao time.Time
	CentralID       int
	MicroID         int
}

func (dbc *DBController) GetCentralMicroByMicroID(microID string) (*CentralMicro, error) {
	stmt := `SELECT versao, data_atualizacao, central_id, micro_id FROM central_micro WHERE micro_id = ?`
	row := dbc.DB.QueryRow(stmt, microID)
	var cm CentralMicro
	err := row.Scan(&cm.Versao, &cm.DataAtualizacao, &cm.CentralID, &cm.MicroID)
	if err != nil {
		return nil, err
	}
	return &cm, nil
}

func (dbc *DBController) InsertCentralMicro(cm CentralMicro) error {
	stmt := `INSERT INTO central_micro (versao, data_atualizacao, central_id, micro_id) VALUES (?, ?, ?, ?)`
	_, err := dbc.DB.Exec(stmt, cm.Versao, cm.DataAtualizacao, cm.CentralID, cm.MicroID)
	return err
}

// ---------- Tabela historico_atualizacao ----------
type HistoricoAtualizacao struct {
	ID              int
	MicroID         int
	DataAtualizacao time.Time
	Versao          string
	Log             string
}

func (dbc *DBController) InsertHistorico(h HistoricoAtualizacao) (int64, error) {
	stmt := `INSERT INTO historico_atualizacao (micro_id, data_atualizacao, versao, log) VALUES (?, ?, ?, ?)`
	res, err := dbc.DB.Exec(stmt, h.MicroID, h.DataAtualizacao, h.Versao, h.Log)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// ---------- Função para gravar payload de treino ----------
func (dbc *DBController) SaveEdgeData(data Treino) error {

	var err error
	_, err = dbc.InsertTreino(data)
	return err
}
