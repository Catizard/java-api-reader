package com.catizard.todo_backend.controller;

import com.baomidou.mybatisplus.core.conditions.query.QueryWrapper;
import com.catizard.todo_backend.common.R;
import com.catizard.todo_backend.common.TokenUtils;
import com.catizard.todo_backend.entity.UserEntity;
import com.catizard.todo_backend.services.TokenService;
import com.catizard.todo_backend.services.UserService;
import com.catizard.todo_backend.vo.LoginUserVo;
import com.catizard.todo_backend.vo.RegisterUserVo;
import com.catizard.todo_backend.vo.UpdatePasswordVo;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpSession;

@RestController
@RequestMapping("user")
public class UserController {

    @Autowired
    private UserService userService;

    @Autowired
    private TokenService tokenService;

    // ----------     info     ----------
    @GetMapping("/info")
    public R getUserInfo(HttpSession session) {
        Long userId = resolveToken(session);
        if (userId == null) {
            return R.error(401, "unauthorized");
        }
        return R.ok().put("data", userService.selectUserById(userId));
    }
    @GetMapping("/info/{userId}")
    public R getUserInfoById(@PathVariable("userId") Long userId) {
        return R.ok().put("data", userService.selectUserById(userId));
    }

    @PostMapping("/password/update")
    public R updatePassword(@RequestBody UpdatePasswordVo updatePasswordVo) {
        //TODO 整合验证码功能?或者单独提供另外一个接口
        return R.ok();
    }


    // ----------     login     ----------
    @PostMapping("/login")
    public R login(@RequestBody LoginUserVo vo, HttpSession session) {
        UserEntity ret = userService.getBaseMapper().selectOne(
                new QueryWrapper<UserEntity>()
                        .eq("user_name", vo.getUserName())
                        .eq("user_password", vo.getUserPassword()));
        if (ret == null) {
            //TODO handle error message
            return R.error(400, "not registered");
        } else {
            String jwt = tokenService.setNewToken(ret.getUserId());
            session.setAttribute("token", jwt);
            return R.ok();
        }
    }

    @GetMapping("/logout")
    public R logout(HttpSession session) {
        tokenService.unsetToken(((String) session.getAttribute("token")));
        return R.ok();
    }

    // ----------     register     ----------
    @GetMapping("/register/sendcode")
    public R sendRegisterCode(@RequestParam("phone") String phone) {
        if (userService.sendRegisterCodeToPhone(phone)) {
            return R.ok();
        }
        return R.error(401, "send failed");
    }

    @PostMapping("/register/register")
    public R register(@RequestBody RegisterUserVo vo) {
        //TODO check if info illegal
        //TODO add cache support
        if (userService.getOne(new QueryWrapper<UserEntity>().eq("user_name", vo.getUserName())) != null) {
            return R.error(401, "name duplicated");
        }
        String userPhone = vo.getUserPhone();
        String checkCode = vo.getCheckCode();
        if (!userService.checkRegisterCode(userPhone, checkCode)) {
            return R.error(401, "check code error");
        }
        //add user to db
        UserEntity userEntity = new UserEntity();
        BeanUtils.copyProperties(vo, userEntity);
        userService.save(userEntity);
        return R.ok();
    }

    private Long resolveToken(HttpSession session) {
        String token = (String) session.getAttribute("token");
        if (token == null) {
            return null;
        }
        return tokenService.verifyToken(token);
    }
}

