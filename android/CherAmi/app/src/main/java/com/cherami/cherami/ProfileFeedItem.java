package com.cherami.cherami;

/**
 * Created by Geoff on 11/28/2014.
 */
public class ProfileFeedItem {
        public String message;
        public String title;
        public String date;

        public ProfileFeedItem(String title, String message, String date) {
            this.title = title;
            this.message = message;
            this.date = date;
        }
}
