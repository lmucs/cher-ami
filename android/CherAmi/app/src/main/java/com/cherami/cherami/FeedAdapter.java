package com.cherami.cherami;

import android.app.Activity;
import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.ImageView;
import android.widget.TextView;

import com.squareup.picasso.Picasso;

import org.json.JSONException;

import java.io.IOException;
import java.net.MalformedURLException;
import java.net.URL;
import java.net.URLConnection;

public class FeedAdapter extends ArrayAdapter<FeedItem> {

    Context context;
    int layoutResourceId;
    FeedItem data[] = null;

    public FeedAdapter(Context context, int layoutResourceId, FeedItem[] data) {
        super(context, layoutResourceId, data);
        this.layoutResourceId = layoutResourceId;
        this.context = context;
        this.data = data;
    }

    public boolean imageUrlChecker(String s){
        URL url;
        try{
            String possibleUrl = s.substring(s.lastIndexOf(' ')+1);
            url = new URL(possibleUrl);
        } catch(MalformedURLException e) {
            return false;
        }
        return true;
    }

    public String getImageUrl(String s){
        return s.substring(s.lastIndexOf(' ')+1);
    }

    public String getContentWithoutImageUrl(String s){
        return s.substring(0,s.lastIndexOf(' '));
    }

    public String processDate(String date){
        return date.substring(0, date.lastIndexOf("T"));
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        View row = convertView;
        FeedHolder holder = null;

        if (row == null) {
            LayoutInflater inflater = ((Activity) context).getLayoutInflater();
            row = inflater.inflate(layoutResourceId, parent, false);

            holder = new FeedHolder();
            holder.txtOwner = (TextView) row.findViewById(R.id.txtOwner);
            holder.txtContent = (TextView) row.findViewById(R.id.txtContent);
            holder.txtDate = (TextView) row.findViewById(R.id.txtDate);
            holder.imgLoad = (ImageView) row.findViewById(R.id.imgLoad);

            row.setTag(holder);
        } else {
            holder = (FeedHolder) row.getTag();
        }

        FeedItem feedItem = data[position];
        try {
            String content = feedItem.msg.getString("content");
            holder.txtOwner.setText(feedItem.msg.getString("author"));
            if(imageUrlChecker(content)) {
                Picasso.with(context).load(getImageUrl(content)).into(holder.imgLoad);
                holder.txtContent.setText(getContentWithoutImageUrl(content));
            } else {
                holder.txtContent.setText(content);
            }
            holder.txtDate.setText(processDate(feedItem.msg.getString("created")));
        } catch (JSONException e){
            e.printStackTrace();
        }

        return row;
    }

    static class FeedHolder {
        TextView txtOwner;
        TextView txtContent;
        TextView txtDate;
        ImageView imgLoad;
    }
}