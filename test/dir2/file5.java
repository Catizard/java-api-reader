ackage com.catizard.todo_backend.dao;

import com.baomidou.mybatisplus.core.mapper.BaseMapper;
import com.catizard.todo_backend.entity.TaskEntity;
import org.apache.ibatis.annotations.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface TaskDao extends BaseMapper<TaskEntity> {
    void addStepCount(@Param("taskId") Long taskId);

    void updateStepComplete(@Param("taskId") Long taskId, @Param("val") int val);

    void decreaseStepCount(@Param("taskId") Long taskId);
}