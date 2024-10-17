package survey

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type SurveyRepository struct {
	Conn *pgx.Conn
}

func NewSurveyRepository(conn *pgx.Conn) *SurveyRepository {
	return &SurveyRepository{
		Conn: conn,
	}
}

func (repo *SurveyRepository) Store(ctx context.Context, input SurveyInput) (*Survey, error) {
	query := "INSERT INTO surveys (user_id, name, phone_no, gender, dob) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var survey Survey
	err := repo.Conn.QueryRow(ctx, query, input.UserId, input.Name, input.PhoneNo, input.Gender, input.Dob).Scan(&survey.Id)
	if err != nil {
		return nil, err
	}
	return repo.getById(ctx, survey.Id)
}

func (repo *SurveyRepository) getById(ctx context.Context, id int) (*Survey, error) {
	var survey Survey
	query := "SELECT * from surveys WHERE id=$1"
	err := repo.Conn.QueryRow(ctx, query, id).
		Scan(&survey.Id, &survey.UserId, &survey.Name, &survey.PhoneNo, &survey.Gender, &survey.Dob, &survey.CreatedAt, &survey.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &survey, nil
}

func (repo *SurveyRepository) getByUserId(ctx context.Context, userId int) ([]Survey, error) {
	query := "SELECT * FROM surveys WHERE user_id=$1"
	rows, _ := repo.Conn.Query(ctx, query, userId)
	surveys, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Survey, error) {
		var survey Survey
		err := row.Scan(&survey.Id, &survey.UserId, &survey.Name, &survey.PhoneNo, &survey.Gender, &survey.Dob, &survey.CreatedAt, &survey.UpdatedAt)
		return survey, err
	})
	if err != nil {
		return nil, err
	}
	return surveys, nil
}

func (repo *SurveyRepository) Update(ctx context.Context, id int, input SurveyInput) (*Survey, error) {
	query := "UPDATE surveys SET name=$1, phone_no=$2, gender=$3, dob=$4 WHERE id=$5"
	res, err := repo.Conn.Exec(ctx, query, input.Name, input.PhoneNo, input.Gender, input.Dob, id)
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, errors.New("No data found")
	}
	return repo.getById(ctx, id)
}
