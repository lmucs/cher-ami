package com.cherami.cherami;

import android.app.DialogFragment;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.os.Bundle;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.MenuItem;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.EditText;
import android.widget.RadioButton;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.apache.http.entity.StringEntity;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.io.UnsupportedEncodingException;
import java.util.Properties;

/**
 * Created by Geoff on 11/20/2014.
 */
public class CreateMessageModal extends DialogFragment{

        SharedPreferences prefs;
        Button createMessageButton;
        Button dismissModalButton;
        View root;

        @Override
        public void onCreate(Bundle savedInstanceState) {
            Context context = getActivity().getApplicationContext();
            prefs = context.getSharedPreferences("com.cherami.cherami", Context.MODE_PRIVATE);
            super.onCreate(savedInstanceState);
        }

        @Override
        public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
            View rootView = inflater.inflate(R.layout.fragment_create_message_modal, container, false);
            getDialog().setTitle("Create New Message");
//        getDialog().setCancelable(true);

            createMessageButton = (Button) rootView.findViewById(R.id.createMessageButton);
            createMessageButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    attemptCreateMessage();
                    dismissModal();
                }
            });
            dismissModalButton = (Button) rootView.findViewById(R.id.dismissModalButton);
            dismissModalButton.setOnClickListener(new View.OnClickListener() {
                @Override
                public void onClick(View v) {
                    dismissModal();
                }
            });
            root = rootView;
            return rootView;
        }

        public String getLocalUrlForApi () {
            AssetManager assetManager = getResources().getAssets();
            InputStream inputStream = null;
            try {
                inputStream = assetManager.open("config.properties");
            } catch (IOException e) {
                e.printStackTrace();
            }
            Properties properties = new Properties();
            try {
                properties.load(inputStream);
            } catch (IOException e) {
                e.printStackTrace();
            }
            return properties.getProperty("myUrl");
        }

        public JSONObject getMessageObjectRequestAsJson () {
            JSONObject jsonParams = new JSONObject();
            EditText messageContent = (EditText) root.findViewById(R.id.messageContent);
            try {
                jsonParams.put("Content", messageContent.getText().toString());
            } catch (JSONException j) {
                System.out.println(j);
            }
            return jsonParams;
        }

        public StringEntity convertJsonUserToStringEntity (JSONObject jsonParams) {
            StringEntity entity = null;
            try {
                entity = new StringEntity(jsonParams.toString());
            } catch (UnsupportedEncodingException i) {
                System.out.println("DONT LIKE TO STRING!");
            }
            return entity;
        }

        public void dismissModal () {
            this.dismiss();
        }

    public void attemptCreateMessage() {
        EditText messageContent = (EditText) root.findViewById(R.id.messageContent);
        Bundle args = new Bundle();
        args.putString("messageValue", messageContent.getText().toString());
        CircleForMessageModal newFragment = new CircleForMessageModal();
        newFragment.setArguments(args);
        newFragment.show(getFragmentManager(), "dialog");


    }



}
