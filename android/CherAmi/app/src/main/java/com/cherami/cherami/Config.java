package com.cherami.cherami;

import java.util.Properties;

public class Config {
    final String directoriesString = "../../../../../../../../../";
    Properties configFile;
    public Config() {
        configFile = new java.util.Properties();
        try {
            configFile.load(this.getClass().getClassLoader().
                    getResourceAsStream(this.directoriesString + "config.cfg"));
        }catch(Exception eta) {
            eta.printStackTrace();
        }
    }

    public String getProperty(String key) {
        String value = this.configFile.getProperty(key);
        return value;
    }
}
