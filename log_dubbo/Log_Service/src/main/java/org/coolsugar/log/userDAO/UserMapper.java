package org.coolsugar.log.userDAO;


import org.springframework.stereotype.Repository;

@Repository
public interface UserMapper {

    UserEntity selectUserById(int id);

    void insertSelective(UserEntity userEntity);

    void updateByPrimaryKeySelective(UserEntity userEntity);
}
