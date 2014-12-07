package com.cherami.cherami;

import android.app.Activity;
import android.app.AlertDialog;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.DialogInterface;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.net.Uri;
import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;
import com.loopj.android.http.RequestParams;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

public class Profile extends Fragment {
    private ListView messageList;
    TextView textElement;
    Context context;
    ProgressDialog dialog;

    public Profile() {

    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        this.context = getActivity().getApplicationContext();
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        String username = ApiHelper.getUsername(context);
        View rootView = inflater.inflate(R.layout.fragment_profile, container, false);
        getProfileFeed(rootView);
        textElement = (TextView) rootView.findViewById(R.id.profileHandle);
        textElement.setText(username);

        return rootView;
    }

    public void getProfileFeed(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);

        client.addHeader("Authorization", token);
        client.get(context,
                   ApiHelper.getLocalUrlForApi(getResources()) + "messages",
                   new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                dialog = ProgressDialog.show(getActivity(), "",
                        "Loading. Please wait...", true);
            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;

                try {
                    responseText = new JSONObject(new String(responseBody)).getString("objects");
                    JSONArray y = new JSONArray(responseText);
                    ProfileFeedItem message_data[] = new ProfileFeedItem[y.length()];

                    for (int x = 0; x < y.length(); x++){
                        message_data[x] = new ProfileFeedItem(new JSONObject(y.get((y.length()-1)-x).toString()).getString("author"),new JSONObject(y.get((y.length()-1)-x).toString()).getString("content"), processDate(new JSONObject(y.get((y.length()-1)-x).toString()).getString("created")));
                    }

                    ProfileFeedAdapter adapter = new ProfileFeedAdapter(getActivity(),
                            R.layout.profile_feed_row, message_data);

                    messageList = (ListView) view.findViewById(R.id.profileFeed);
                    messageList.setAdapter(adapter);

                } catch (JSONException j) {

                }
                dialog.dismiss();

            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                dialog.dismiss();
                String responseText = null;

                try {
                    if (!NetworkCheck.isConnected(errorResponse)) {
                        new AlertDialog.Builder(getActivity())
                                .setTitle("Network Error")
                                .setMessage("You're not connected to the network :(")
                                .setNegativeButton(android.R.string.ok, new DialogInterface.OnClickListener() {
                                    public void onClick(DialogInterface dialog, int which) {
                                        // do nothing
                                    }
                                })
                                .setIcon(android.R.drawable.ic_dialog_alert)
                                .show();

                    } else {
                        responseText = new JSONObject(new String(errorResponse)).getString("reason");
                        Toast toast = Toast.makeText(getActivity().getApplicationContext(), responseText, Toast.LENGTH_LONG);
                        toast.show();
                    }
                } catch (JSONException j) {

                }

            }
        });
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    public void showLogin(View view) {
        Intent intent = new Intent(getActivity(), LoginActivity.class);
        startActivity(intent);
    }
}
