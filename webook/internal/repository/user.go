package repository

import (
	"context"
	"database/sql"

	"github.com/spigcoder/LittleBook/webook/internal/domain"
	"github.com/spigcoder/LittleBook/webook/internal/repository/cache"
	"github.com/spigcoder/LittleBook/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDao, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao: dao,
		cache: cache,
	}
}

func (repo *UserRepository) FindByPhone(c context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(c, phone)
	if err!= nil {
		return domain.User{}, err
	}	
	return repo.TransformDaoToDomain(u), nil
}

func (repo *UserRepository) FindByEmail(c context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(c, email)
	if err != nil {
		return domain.User{}, err
	}

	return repo.TransformDaoToDomain(u), nil 
}

func (repo *UserRepository) Create(c context.Context, u domain.User) error {
	err := repo.dao.Insert(c, repo.TranformDomainToDao(u))
	//这里可以做缓存，这里是注册，所以缓存的key可以使用邮箱，这样再次登录可以直接从缓存中获取用户信息
	//但是登录没必要做缓存
	return err
}

func (repo *UserRepository) Edit(c context.Context, u domain.User) error {
	//这里涉及缓存一致性问题，可以考虑先删除缓存，然后再更新数据库
	err := repo.dao.Edit(c, dao.User{
		Birthday: u.Birthday,
		Intro:    u.Intro,
		UserName: u.UserName,
		Id:       u.Id,
	})
	return err
}

func (repo *UserRepository) GetUserById(c context.Context, id int64) (user domain.User, err error) {
	// 先从缓存中获取用户信息
	user, err = repo.cache.GetById(c, id)
	if err == cache.KeyNotExist {
		// 缓存中没有用户信息，从数据库中获取
		daoUser, err := repo.dao.FindById(c, id)	
		// 这里有问题，如果说有人大量访问db中没有的问题，可能会导致缓存穿透问题，可以辅助使用布隆过滤器
		if err != nil {
			return domain.User{}, err
		}
		user := repo.TransformDaoToDomain(daoUser)
		err = repo.cache.Set(c, user)
		if err!= nil {
			//TODO 这里要记录日志
		}
		return user, nil
	}
	// 这里又有问题了，缓存如果出错怎么办
	//如果这里继续查询缓存肯能会导致缓存雪崩问题（缓存崩溃），那这时如果还要查询数据库就一定要保护好数据库
	//可以使用限流方法来进行处理
	//也可以不查询，直接返回空，但是可能会导致用户体验不好
	if err!= nil {
		//保守做法，返回空
		return domain.User{}, err
	}
	return user, nil
}


func (repo *UserRepository) TranformDomainToDao(user domain.User) dao.User {
	return dao.User{
		Id:       user.Id,
		Email:    sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Password: user.Password,
		Birthday: user.Birthday,
		UserName: user.UserName,
		Intro:    user.Intro,
	}
}

func (repo *UserRepository) TransformDaoToDomain(daoUser dao.User) domain.User {
	return domain.User{
		Id:       daoUser.Id,
		Email:    daoUser.Email.String,
		Phone:    daoUser.Phone.String,
		Password: daoUser.Password,
		Birthday: daoUser.Birthday,
		UserName: daoUser.UserName,
		Intro:    daoUser.Intro,
	}
}
