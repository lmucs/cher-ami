package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
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


/**
 * A simple {@link Fragment} subclass.
 * Activities that contain this fragment must implement the
 * {@link Profile.OnFragmentInteractionListener} interface
 * to handle interaction events.
 * Use the {@link Profile#newInstance} factory method to
 * create an instance of this fragment.
 *
 */
public class Profile extends Fragment {
    // TODO: Rename parameter arguments, choose names that match
    // the fragment initialization parameters, e.g. ARG_ITEM_NUMBER
    private static final String ARG_PARAM1 = "param1";
    private static final String ARG_PARAM2 = "param2";

    // TODO: Rename and change types of parameters
    private String mParam1;
    private String mParam2;
    private ListView messageList;
    SharedPreferences prefs;

    TextView textElement;
    private OnFragmentInteractionListener mListener;

    /**
     * Use this factory method to create a new instance of
     * this fragment using the provided parameters.
     *
     * @param param1 Parameter 1.
     * @param param2 Parameter 2.
     * @return A new instance of fragment Profile.
     */
    // TODO: Rename and change types and number of parameters
    public static Profile newInstance(String param1, String param2) {
        Profile fragment = new Profile();
        Bundle args = new Bundle();
        args.putString(ARG_PARAM1, param1);
        args.putString(ARG_PARAM2, param2);
        fragment.setArguments(args);
        return fragment;
    }
    public Profile() {
        // Required empty public constructor
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        Context context = getActivity().getApplicationContext();
        prefs = context.getSharedPreferences(
                "com.cherami.cherami", Context.MODE_PRIVATE);
        super.onCreate(savedInstanceState);
        if (getArguments() != null) {
            mParam1 = getArguments().getString(ARG_PARAM1);
            mParam2 = getArguments().getString(ARG_PARAM2);
        }
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        String userKey = "com.cherami.cherami.username";
        String username = prefs.getString(userKey, null);
        View rootView = inflater.inflate(R.layout.fragment_profile, container, false);
        getProfileFeed(rootView);
        textElement = (TextView) rootView.findViewById(R.id.profileHandle);
        textElement.setText(username);
        // Inflate the layout for this fragment
        return rootView;
    }

    public void getProfileFeed(final View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String token = ApiHelper.getSessionToken(prefs);

        client.addHeader("Authorization", token);
        client.get(getActivity().getApplicationContext(),
                   ApiHelper.getLocalUrlForApi(getResources()) + "messages",
                   new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {

            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(responseBody)).getString("objects");
                    System.out.println(responseText);
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

            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                // called when response HTTP status is "4XX" (eg. 401, 403, 404)

                String responseText = null;
                try {
                    responseText = new JSONObject(new String(errorResponse)).getString("reason");

                } catch (JSONException j) {
                    System.out.println(j);
                }

            }
        });
    }

    // TODO: Rename method, update argument and hook method into UI event
    public void onButtonPressed(Uri uri) {
        if (mListener != null) {
            mListener.onFragmentInteraction(uri);
        }
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    @Override
    public void onAttach(Activity activity) {
        super.onAttach(activity);
        try {
            mListener = (OnFragmentInteractionListener) activity;
        } catch (ClassCastException e) {
            throw new ClassCastException(activity.toString()
                    + " must implement OnFragmentInteractionListener");
        }
    }

    @Override
    public void onDetach() {
        super.onDetach();
        mListener = null;
    }

    public void showLogin(View view) {
        Intent intent = new Intent(getActivity(), LoginActivity.class);
        startActivity(intent);
    }

    public interface OnFragmentInteractionListener {
        // TODO: Update argument type and name
        public void onFragmentInteraction(Uri uri);
    }

}
