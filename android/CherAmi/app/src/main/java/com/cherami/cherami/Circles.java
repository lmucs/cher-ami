package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.res.AssetManager;
import android.net.Uri;
import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.ListView;
import android.widget.Spinner;
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
import java.lang.reflect.Array;
import java.util.Properties;


/**
 * A simple {@link Fragment} subclass.
 * Activities that contain this fragment must implement the
 * {@link Circles.OnFragmentInteractionListener} interface
 * to handle interaction events.
 * Use the {@link Circles#newInstance} factory method to
 * create an instance of this fragment.
 *
 */
public class Circles extends Fragment {
    // TODO: Rename parameter arguments, choose names that match
    // the fragment initialization parameters, e.g. ARG_ITEM_NUMBER
    private static final String ARG_PARAM1 = "param1";
    private static final String ARG_PARAM2 = "param2";
    private ListView circleList;
    // TODO: Rename and change types of parameters
    private String mParam1;
    private String mParam2;
    private Spinner spinner;
    SharedPreferences prefs;

    private OnFragmentInteractionListener mListener;

    /**
     * Use this factory method to create a new instance of
     * this fragment using the provided parameters.
     *
     * @param param1 Parameter 1.
     * @param param2 Parameter 2.
     * @return A new instance of fragment Circles.
     */
    // TODO: Rename and change types and number of parameters
    public static Circles newInstance(String param1, String param2) {
        Circles fragment = new Circles();
        Bundle args = new Bundle();
        args.putString(ARG_PARAM1, param1);
        args.putString(ARG_PARAM2, param2);
        fragment.setArguments(args);
        return fragment;
    }
    public Circles() {
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

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {

        View rootView = inflater.inflate(R.layout.fragment_circles, container, false);

        getCircles(rootView);

        // Get the filter value
        spinner = (Spinner) rootView.findViewById(R.id.filter_spinner);
        spinner.setOnItemSelectedListener(new AdapterView.OnItemSelectedListener() {
            public void onItemSelected(AdapterView<?> parent, View view,
                                       int position, long id) {
                String value = spinner.getSelectedItem().toString();
                System.out.println("Value: " + value);
            }

            @Override
            public void onNothingSelected(AdapterView<?> parent) {

            }
        });

        // Inflate the layout for this fragment
        return rootView;
    }

    public void getCircles(View view) {
        AsyncHttpClient client = new AsyncHttpClient();
        String sessionKey = "com.cherami.cherami.sessionid";
        String sessionid = prefs.getString(sessionKey, null);
        String userKey = "com.cherami.cherami.username";
        String username = prefs.getString(userKey, null);
        RequestParams params = new RequestParams();
        params.put("user", username);
        final String[] circleArray = new String[100];
        final View view2 = view;


        client.addHeader("Authorization", sessionid);
        client.get(getActivity().getApplicationContext(), "http://" + getLocalUrlForApi() + "/api/circles", params, new AsyncHttpResponseHandler() {

            @Override
            public void onStart() {
                // called before request is started
                System.out.println("STARTING GET REQUEST");

            }

            @Override
            public void onSuccess(int statusCode, Header[] headers, byte[] responseBody) {
                String responseText = null;
                try {
                    responseText = new JSONObject(new String(responseBody)).getString("Results");
                    JSONArray y = new JSONArray(responseText);
                    for(int x = 0; x < y.length(); x++){
                        circleArray[x] = y.get(x).toString();
                    }
                    Circle circle_data[] = new Circle[100];
                    for (int x = 0; x < circleArray.length; x++){
                        circle_data[x] = new Circle(circleArray[x]);
                    }

                    CircleAdapter adapter = new CircleAdapter(getActivity(),
                            R.layout.circle_item_row, circle_data);


                    circleList = (ListView) view2.findViewById(R.id.circleList);

                    circleList.setAdapter(adapter);
                    String s = new JSONObject(new JSONArray(responseText).get(0).toString()).getString("c.name");
                    System.out.println(s);
                } catch (JSONException j) {
                    System.out.println("Dont like JSON");
                }


            }

            @Override
            public void onFailure(int statusCode, Header[] headers, byte[] errorResponse, Throwable error) {
                // called when response HTTP status is "4XX" (eg. 401, 403, 404)

                String responseText = null;
                try {
                    responseText = new JSONObject(new String(errorResponse)).getString("Reason");

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

    /**
     * This interface must be implemented by activities that contain this
     * fragment to allow an interaction in this fragment to be communicated
     * to the activity and potentially other fragments contained in that
     * activity.
     * <p>
     * See the Android Training lesson <a href=
     * "http://developer.android.com/training/basics/fragments/communicating.html"
     * >Communicating with Other Fragments</a> for more information.
     */
    public interface OnFragmentInteractionListener {
        // TODO: Update argument type and name
        public void onFragmentInteraction(Uri uri);
    }

}
