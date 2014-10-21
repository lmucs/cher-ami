package com.cherami.cherami;

public class NewUser {

    private String handle;
    private String email;
    private String password;
    private String confirmpassword;


    public NewUser(){

    }

    public NewUser(String handle, String email, String password, String confirmpassword){
        this.handle = handle;
        this.email = email;
        this.password = password;
        this.confirmpassword = confirmpassword;
    }

    public void setHandle(String s) {
        handle = s;
    }

    public void setEmail(String s) {
        email = s;
    }

    public void setPassword(String s) {
       password = s;
    }

    public void setConfirmpassword(String s) {
        confirmpassword = s;
    }

    public String getHandle() {
        return this.handle;
    }

    public String getEmail() {
        return this.email;
    }

    public String getPassword() {
        return this.password;
    }

    public String getConfirmpassword() {
        return this.confirmpassword;
    }
}