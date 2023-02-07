$(function(){
    // 產品頁
    var table = $('#productTable').DataTable({
        language: {
            url: "static/vendor/dataTables/languages/ch-tw.json"  
        },
        processing:true,
        serverSide:true,
        ajax: {
            url: "/api/getProductsList",
            type: "GET",
        },
        columns:[
            {
                data: 'name', 
                title:"產品名稱",
            },
            {
                data: 'amount', 
                title:"數量"
            },
            {
                data: 'formatTime',
                title:"更新時間"
            },
            {
                data: 'id', 
                title: "修改存貨",
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<input type="text" data-id="'+data+'" style="margin-right:3px;">';
                    output += '<button type="button" class="btn btn-success btn-sm">修改</button></div>';
                    return output;
                }
            },
        ],
        columnDefs:[
            {
                targets: 3,
                orderable: false,
                responsivePriority: 2,
            },
            {
                targets: 0,
                responsivePriority: 1,
            }
        ],
    });
});