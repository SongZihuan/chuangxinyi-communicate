package service

import (
	errors "github.com/wuntsong-org/wterrors"
	"strconv"

	"gitee.com/wuntsong/chuangxinyi-communicate/cache"
	"gitee.com/wuntsong/chuangxinyi-communicate/dao"
	"gitee.com/wuntsong/chuangxinyi-communicate/logger"
	"gitee.com/wuntsong/chuangxinyi-communicate/model"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils"
	"gitee.com/wuntsong/chuangxinyi-communicate/utils/sqlcnd"
)

var UserScoreService = newUserScoreService()

func newUserScoreService() *userScoreService {
	return &userScoreService{}
}

type userScoreService struct {
}

func (s *userScoreService) Get(id int64) *model.UserScore {
	return dao.UserScoreDao.Get(id)
}

func (s *userScoreService) Take(where ...interface{}) *model.UserScore {
	return dao.UserScoreDao.Take(where...)
}

func (s *userScoreService) Find(cnd *sqlcnd.SqlCnd) []model.UserScore {
	return dao.UserScoreDao.Find(cnd)
}

func (s *userScoreService) FindOne(cnd *sqlcnd.SqlCnd) *model.UserScore {
	return dao.UserScoreDao.FindOne(cnd)
}

func (s *userScoreService) List(cnd *sqlcnd.SqlCnd) (list []model.UserScore, paging *sqlcnd.Paging) {
	return dao.UserScoreDao.List(cnd)
}

func (s *userScoreService) Create(t *model.UserScore) errors.WTError {
	return dao.UserScoreDao.Create(t)
}

func (s *userScoreService) Update(t *model.UserScore) errors.WTError {
	return dao.UserScoreDao.Update(t)
}

func (s *userScoreService) Updates(id int64, columns map[string]interface{}) errors.WTError {
	return dao.UserScoreDao.Updates(id, columns)
}

func (s *userScoreService) UpdateColumn(id int64, name string, value interface{}) errors.WTError {
	return dao.UserScoreDao.UpdateColumn(id, name, value)
}

func (s *userScoreService) Delete(id int64) {
	dao.UserScoreDao.Delete(id)
}

func (s *userScoreService) GetByUserId(userId int64) *model.UserScore {
	return s.FindOne(sqlcnd.NewSqlCnd().Eq("user_id", userId))
}

func (s *userScoreService) CreateOrUpdate(t *model.UserScore) errors.WTError {
	if t.ID > 0 {
		return s.Update(t)
	} else {
		return s.Create(t)
	}
}

// IncrementCreateTopicScore 发帖获积分
func (s *userScoreService) IncrementPostTopicScore(topic *model.Topic) {
	config := SettingService.GetSetting()
	if config.ScoreConfig.PostTopicScore <= 0 {
		logger.Logger.Info("请配置发帖积分")
		return
	}
	err := s.addScore(topic.UserId, config.ScoreConfig.PostTopicScore, model.EntityTypeTopic,
		strconv.FormatInt(topic.ID, 10), "发表话题")
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

// IncrementPostCommentScore 跟帖获积分
func (s *userScoreService) IncrementPostCommentScore(comment *model.Comment) {
	// 非话题跟帖，跳过
	if comment.EntityType != model.EntityTypeTopic {
		return
	}
	config := SettingService.GetSetting()
	if config.ScoreConfig.PostCommentScore <= 0 {
		logger.Logger.Info("请配置跟帖积分")
		return
	}
	err := s.addScore(comment.UserId, config.ScoreConfig.PostCommentScore, model.EntityTypeComment,
		strconv.FormatInt(comment.ID, 10), "发表跟帖")
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

// Increment 增加分数
func (s *userScoreService) Increment(userId int64, score int, sourceType, sourceId, description string) errors.WTError {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, score, sourceType, sourceId, description)
}

// Decrement 减少分数
func (s *userScoreService) Decrement(userId int64, score int, sourceType, sourceId, description string) errors.WTError {
	if score <= 0 {
		return errors.New("分数必须为正数")
	}
	return s.addScore(userId, -score, sourceType, sourceId, description)
}

// addScore 加分数，也可以加负数
func (s *userScoreService) addScore(userId int64, score int, sourceType, sourceId, description string) errors.WTError {
	if score == 0 {
		return errors.New("分数不能为0")
	}
	userScore := s.GetByUserId(userId)
	if userScore == nil {
		userScore = &model.UserScore{
			UserId:     userId,
			CreateTime: utils.NowTimestamp(),
		}
	}
	userScore.Score = userScore.Score + score
	userScore.UpdateTime = utils.NowTimestamp()
	if err := s.CreateOrUpdate(userScore); err != nil {
		return errors.WarpQuick(err)
	}

	scoreType := model.ScoreTypeIncr
	if score < 0 {
		scoreType = model.ScoreTypeDecr
	}
	err := UserScoreLogService.Create(&model.UserScoreLog{
		UserId:      userId,
		SourceType:  sourceType,
		SourceId:    sourceId,
		Description: description,
		Type:        scoreType,
		Score:       score,
		CreateTime:  utils.NowTimestamp(),
	})
	if err == nil {
		cache.UserCache.InvalidateScore(userId)
	}
	return errors.WarpQuick(err)
}
