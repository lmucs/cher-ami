package com.cherami.cherami;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.DialogInterface;
import android.widget.Toast;

import org.json.JSONException;
import org.json.JSONObject;

public class NetworkCheck {
    public static Boolean isConnected (byte[] response) {
        if (response == null) {
            return false;
        } else {
            return true;
        }
    }

    public static void displayNetworkErrorModal (Activity activity) {
        new AlertDialog.Builder(activity)
            .setTitle("Network Error")
            .setMessage("You're not connected to the network :(")
            .setNegativeButton(android.R.string.ok, new DialogInterface.OnClickListener() {
                public void onClick(DialogInterface dialog, int which) {
                    // do nothing
                }
            })
            .setIcon(android.R.drawable.ic_dialog_alert)
            .show();
    }
}
