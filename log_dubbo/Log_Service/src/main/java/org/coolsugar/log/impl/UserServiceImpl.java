package org.coolsugar.log.impl;

import org.coolsugar.log.po.User;
import org.coolsugar.log.userDAO.UserEntity;
import org.coolsugar.log.userDAO.UserMapper;
import org.coolsugar.log.service.UserService; //api层
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service("userService")
public class UserServiceImpl implements UserService {

    @Autowired
    private UserMapper userMapper;

    public User getUserById(int id) {
        UserEntity userEntity = userMapper.selectUserById(id);

        User user = new User(); //test-api层

        user.setId(userEntity.getId());
        user.setCreator(userEntity.getCreator());
        user.setGmtCreate(userEntity.getGmtCreate());
        user.setModifier(userEntity.getModifier());
        user.setGmtModified(userEntity.getGmtModified());
        user.setIsDeleted(userEntity.getIsDeleted());
        user.setName(userEntity.getName());
        user.setUserType(userEntity.getUserType());
        user.setUserId(userEntity.getUserId());
        user.setAppKey(userEntity.getAppKey());

        return user;
    }

    public void saveUser(User user) {
        Integer id = user.getId();

        UserEntity userEntity = new UserEntity();
        userEntity.setId(user.getId());
        userEntity.setCreator(user.getCreator());
        userEntity.setName(user.getName());

        if(id == null) { //
            //userMapper.insertSelective(userEntity);
            System.out.println("插入成功！");
        }
        else {
            //userMapper.updateByPrimaryKeySelective(userEntity);
            System.out.println("更新成功！");
        }
    }

    public void sayHi() {
        System.out.println("hihihi...");
    }
}
