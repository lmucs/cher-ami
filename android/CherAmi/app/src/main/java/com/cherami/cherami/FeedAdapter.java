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

import java.net.MalformedURLException;
import java.net.URL;

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
            holder.txtOwner.setText(feedItem.msg.getString("author"));
            holder.txtContent.setText(feedItem.msg.getString("content"));
            holder.txtDate.setText(feedItem.msg.getString("created"));
        } catch (JSONException e){
            e.printStackTrace();
        }
        if(!(feedItem.img.equals(""))){
            Picasso.with(context).load(feedItem.img).into(holder.imgLoad);
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