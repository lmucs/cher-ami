package com.cherami.cherami;

import android.content.res.AssetManager;
import android.content.res.Resources;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

/**
 * Created by crashprophet on 11/30/14.
 */
public class ApiHelper {

    public static String getLocalUrlForApi (Resources resources) {
        AssetManager assetManager = resources.getAssets();
        Properties properties = new Properties();
        try {
            properties.load(assetManager.open("config.properties"));
        } catch (IOException e) {
            e.printStackTrace();
        }

        return properties.getProperty("myUrl");
    }
}
