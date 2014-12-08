package com.cherami.cherami;

import android.app.ProgressDialog;
import android.content.Context;
import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ListView;
import android.widget.Spinner;
import android.widget.Toast;

import com.loopj.android.http.AsyncHttpClient;
import com.loopj.android.http.AsyncHttpResponseHandler;

import org.apache.http.Header;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

public class Feed extends Fragment {
    private ListView feedList;
    private Spinner spinner;
    Context context;
    FeedAdapter adapter;
    ProgressDialog dialog;

    public Feed() {

    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        this.context = getActivity().getApplicationContext();
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        View rootView = inflater.inflate(R.layout.fragment_feed, container, false);
        getFeed(rootView);

        // Get the filter value
        spinner = (Spinner) rootView.findViewById(R.id.filter_spinner);
        spinner.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            public void onItemSelected(AdapterView<?> parent, View view,
                                       int position, long id) {
                String value = spinner.getSelectedItem().toString();
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {

            }
        });

        // Inflate the layout for this fragment
        return rootView;
    }

    public void setFeedAdapter (FeedAdapter feedAdapter) {
        this.adapter = feedAdapter;
    }

    public void getFeed(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(context);

        client.addHeader("Authorization", token);
        client.get(context, ApiHelper.getLocalUrlForApi(getResources()) + "messages",
                new AsyncHttpResponseHandler() {

                    @Override
                    public void onStart() {
                        dialog = ProgressDialog.show(getActivity(), "",
                                "Loading. Please wait...", true);
                    }

                    @Override
                    public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                        JSONArray responseText;

                        try {
                            responseText = new JSONArray(new String(responseBody));
                            FeedItem feed_data[] = new FeedItem[responseText.length()];

                            for (int x = 0; x < responseText.length(); x++) {
                                feed_data[x] = new FeedItem(new JSONObject(responseText.get(x).toString()));
                            }

                            final FeedAdapter adapter = new FeedAdapter(getActivity(),
                                    R.layout.feed_item_row, feed_data);
                            Feed.this.setFeedAdapter(adapter);
                            feedList = (ListView) view.findViewById(R.id.feedList);
                            feedList.setAdapter(adapter);

                        } catch (JSONException j) {

                        }
                        dialog.dismiss();
                    }

                    @Override
                    public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                        dialog.dismiss();
                        String responseText = null;

                        try {
                            if (!ErrorHandle.isNetworkConnected(errorResponse)) {
                                ErrorHandle.displayNetworkErrorModal(getActivity());

                            } else {
                                responseText = new JSONObject(new String(errorResponse)).getString("reason");
                                if (ErrorHandle.isTokenExpired(responseText)) {
                                    ErrorHandle.displayTokenModal(getActivity());
                                }
                            }
                        } catch (JSONException j) {

                        }

                    }
                });
    }
}
