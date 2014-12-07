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
            holder.txtTitle = (TextView) row.findViewById(R.id.txtTitle);
            holder.imgLoad = (ImageView) row.findViewById(R.id.imgLoad);

            row.setTag(holder);
        } else {
            holder = (FeedHolder) row.getTag();
        }

        FeedItem feedItem = data[position];
        holder.txtTitle.setText(feedItem.title);
        if(!(feedItem.img.equals(""))){
            Picasso.with(context).load(feedItem.img).into(holder.imgLoad);
        }

        return row;
    }

    static class FeedHolder {
        TextView txtTitle;
        ImageView imgLoad;
    }
}