package com.cherami.cherami;

/**
 * Created by goalsman on 12/7/14.
 */
public class NetworkCheck {
    public static Boolean isConnected (byte[] response) {
        if (response == null) {
            return false;
        } else {
            return true;
        }
    }
}
