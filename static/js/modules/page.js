/**
 * Copyright (c) 2018 phachon@163.com
 */

var Page = {

    /**
     * ajax Save
     * @param element
     * @returns {boolean}
     */
    ajaxSave: function (element, sendEmail, isAutoFollow) {

        /**
         * 成功信息条
         * @param message
         * @param data
         */
        function successBox(message, data) {
            Layers.successMsg(message)
        }

        /**
         * 失败信息条
         * @param message
         * @param data
         */
        function failed(message, data) {
            Layers.failedMsg(message)
        }

        /**
         * request success
         * @param result
         */
        function response(result) {
            //console.log(result)
            if (result.code == 0) {
                failed(result.message, result.data);
            }
            if (result.code == 1) {
                // remove storage
                var documentId = $("input[name='document_id']").val();
                var storageId = "mm_wiki_doc_"+documentId;
                Storage.remove(storageId);
                successBox(result.message, result.data);
            }
            if (result.redirect.url) {
                var sleepTime = result.redirect.sleep || 3000;
                setTimeout(function () {
                    parent.location.href = result.redirect.url;
                }, sleepTime);
            }
        }

        var containerHtml = '<div class="container-fluid" style="padding: 20px 20px 0 20px">';
        containerHtml += '<div style="margin-top: 8px;text-align: left">';

        if (isAutoFollow !== "1") {
            containerHtml += '<label>&nbsp;&nbsp;关注该文档&nbsp;</label><input type="checkbox" name="is_follow_doc" checked="checked" value="1">';
        }

        containerHtml += "</div>";

        layer.open({
            title: '<i class="fa fa-volume-up"></i> 请确认修改',
            type: 1,
            area: ['380px', '232px'],
            content: containerHtml,
            btn: ['确定','取消'],
            yes: function(index, layero){
                var isFollowDoc = "0";
                if (isAutoFollow !== "1") {
                    var followDocCheck = $("input[type='checkbox'][name='is_follow_doc']").is(':checked');
                    if (followDocCheck) {
                        isFollowDoc = "1";
                    }
                }
                layer.close(index);
                var options = {
                    dataType: 'json',
                    success: response,
                    data: {'is_follow_doc': isFollowDoc}
                };
                $(element).ajaxSubmit(options);
            },
            btn2: function(index, layero){
                layer.close(index);
            }
        });

        return false;
    },

    /**
     * cancel save
     */
    cancelSave: function (title, url) {
        title = '<i class="fa fa-volume-up"></i> '+title;
        layer.confirm(title, {
            btn: ['是','否'],
            skin: Layers.skin,
            btnAlign: 'c',
            title: "<i class='fa fa-warning'></i><strong> 警告</strong>"
        }, function() {
            var documentId = $("input[name='document_id']").val();
            var storageId = "mm_wiki_doc_"+documentId;
            Storage.remove(storageId);
            parent.location = url
        }, function() {

        });
    },

    /**
     * upload attachment
     */
    attachment: function (documentId) {
        layer.open({
            type: 2,
            skin: Layers.skin,
            title: '<strong>附件</strong>',
            shadeClose: true,
            shade : 0.1,
            resize: false,
            maxmin: false,
            area: ["900px", "500px"],
            content: "/attachment/page?document_id="+documentId,
            padding:"10px"
        });
    },

    /**
     * 错误提示
     * @param element
     * @param message
     */
    uploadErrorBox: function (element, message) {
        $(element).html('');
        $(element).removeClass('hide');
        $(element).addClass('alert alert-danger');
        $(element).append('<a class="close" href="#" onclick="$(this).parent().hide();">×</a>');
        $(element).append('<strong><i class="glyphicon glyphicon-remove-circle"></i> 上传失败：</strong>');
        $(element).append(message);
        $(element).show();
    },
};