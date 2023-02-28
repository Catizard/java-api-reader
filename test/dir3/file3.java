package com.catizard.todo_backend.services.impl;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.catizard.todo_backend.common.DeptUtils;
import com.catizard.todo_backend.dao.DeptRelationDao;
import com.catizard.todo_backend.entity.DeptEntity;
import com.catizard.todo_backend.entity.DeptRelationEntity;
import com.catizard.todo_backend.entity.UserEntity;
import com.catizard.todo_backend.services.DeptRelationService;
import com.catizard.todo_backend.services.DeptService;
import com.catizard.todo_backend.services.UserService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

@Service
@Slf4j
public class DeptRelationServiceImpl extends ServiceImpl<DeptRelationDao, DeptRelationEntity> implements DeptRelationService {
    @Autowired
    private UserService userService;
    @Autowired
    private DeptService deptService;
    @Override
    public boolean join(Long deptId, Long userId, Integer role) {
        DeptRelationEntity relation = new DeptRelationEntity();
        relation.setDeptId(deptId);
        relation.setUserId(userId);
        relation.setUserRole(role);
        return super.save(relation);
    }

    @Override
    public boolean join(Long deptId, Long userId) {
        return this.join(deptId, userId, DeptUtils.DEPT_USER_ROLE_MEMBER);
    }

    @Override
    public List<DeptRelationEntity> listByUserId(Long userId) {
        return baseMapper.selectList(new QueryWrapper<DeptRelationEntity>().eq("user_id", userId));
    }

    @Override
    public List<UserEntity> listUserInDept(Long deptId) {
        List<DeptRelationEntity> deptRelationEntities = baseMapper.selectList(new QueryWrapper<DeptRelationEntity>().eq("dept_id", deptId));
        List<Long> userIds = deptRelationEntities.stream().map(DeptRelationEntity::getUserId).collect(Collectors.toList());
        return userService.selectBatchByUserIds(userIds);
    }

    @Override
    public List<DeptEntity> listDeptWithUser(Long userId) {
        List<DeptRelationEntity> deptRelationEntities = listByUserId(userId);
        List<Long> deptIds = deptRelationEntities.stream().map(DeptRelationEntity::getDeptId).collect(Collectors.toList());
        // 如果给listByIds方法传入一个空列表，拼接sql时会出现 in () 导致错误
        if (deptIds.isEmpty()) {
            return null;
        }
        List<DeptEntity> deptEntities = (List<DeptEntity>) deptService.listByIds(deptIds);
        log.debug("list id={}, depts={}", userId, deptEntities);
        return deptEntities;
    }
}
