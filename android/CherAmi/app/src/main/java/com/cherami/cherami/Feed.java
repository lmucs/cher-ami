package com.cherami.cherami;

import android.app.Activity;
import android.net.Uri;
import android.os.Bundle;
import android.app.Fragment;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.AdapterView;
import android.widget.ListView;
import android.widget.Spinner;

import java.net.MalformedURLException;
import java.net.URL;

public class Feed extends Fragment {
    private ListView feedList;
    private Spinner spinner;

    public Feed() {

    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
    }

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {
        // Inflate the layout for this fragment
        View rootView = inflater.inflate(R.layout.fragment_feed, container, false);

        FeedItem feed_data[] = new FeedItem[]
                {
                        new FeedItem("This is an image posted by Willy Hugestud","http://www.amovieaweek.com/images/social1.png"),
                        new FeedItem("Here's some text posted by ThatHalfKorean","http://s2.quickmeme.com/img/25/25418a62177a75a6b412ffcc6328ddd37bc8495df2442fc62f1b1fa5408b65a1.jpg"),
                        new FeedItem("I just took a 25 minutes bathroom break posted by CrashProphet")
                };

        FeedAdapter adapter = new FeedAdapter(this.getActivity(),
                R.layout.feed_item_row, feed_data);


        feedList = (ListView)rootView.findViewById(R.id.feedList);

        feedList.setAdapter(adapter);

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
}
